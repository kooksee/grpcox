package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fullstorydev/grpchan/httpgrpc"
	"github.com/gusaul/grpcox/internal/proto/demov1pb"
	"github.com/gusaul/grpcox/web/ace"
	"github.com/gusaul/grpcox/web/jsutil"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pubgo/xerror"
	"honnef.co/go/js/dom/v2"
	"io/ioutil"
	"net/http"
	"net/url"
)

var doc = dom.GetWindow().Document()

func (h *Home) OnInit() {
	fmt.Println("OnInit", h.Mounted())
	h.data = make(map[string]*Request)
	h.tables = make(map[string]bool)
	h.picker = (&jsutil.FilePicker{ID: "hiddenFilePicker", Multiple: false}).Accept("image/*")
	h.tableHidden = true
	h.req = new(Request)

	u, err := url.Parse(fmt.Sprintf("http://127.0.0.1:6969/grpc"))
	if err != nil {
		panic(err)
	}
	h.cc = demov1pb.NewTransportClient(&httpgrpc.Channel{
		Transport: http.DefaultTransport,
		BaseURL:   u,
	})
}

func (h *Home) OnDismount() {
	fmt.Println("OnDismount", h.Mounted())
}

func (h *Home) OnMount(ctx app.Context) {
	h.target = "localhost:50051"
	fmt.Println("OnMount", h.Mounted())

	fmt.Println("OnMount", h.data)

	fmt.Println(jsutil.LoadJs("https://ace.c9.io/build/src/ace.js"))

	ctx.Async(func() {
		rsp, err := http.Get("/api/requests")
		xerror.Panic(err)

		var req map[string][]*Request
		dtBytes, err := ioutil.ReadAll(rsp.Body)
		xerror.Panic(err)
		xerror.Panic(json.Unmarshal(dtBytes, &req))

		for _, r := range req["data"] {
			h.data[r.ID] = r
		}
		ctx.Dispatch(func(context app.Context) {})
	})

	ctx.Async(func() {
		var ret, _ = h.cc.Unary(context.Background(), &demov1pb.Message{Hello: "hello"})
		fmt.Println(ret)
	})

	//ctx.Async(func() {
	//	defer xerror.RecoverAndExit()
	//	var ret, err = h.cc.ServerStream(context.Background(), &demov1pb.Message{
	//		Hello: "world",
	//	})
	//	if err != nil {
	//		return
	//	}
	//
	//	for {
	//		mm, err := ret.Recv()
	//		xerror.Panic(err)
	//		if err != nil {
	//			break
	//		} else {
	//			fmt.Println("ServerStream value:", mm)
	//		}
	//	}
	//})
}

func (h *Home) OnNav(ctx app.Context) {
	fmt.Println("OnNav", h.Mounted())
	p := ctx.Page()
	p.SetTitle("Image to Factorio Blueprint")
	p.SetAuthor("mlctrez")
	p.SetKeywords("factorio, blueprint, image, convert")
	p.SetDescription("A progressive web application for converting images to factorio tile blueprints.")
}

func page(uis ...app.UI) app.UI {
	return app.Div().Class("pf-c-page").Body(uis...)
}

func pageSidebar(uis ...app.UI) app.UI {
	return app.Div().Class("pf-c-page__sidebar").Body(
		app.Div().Class("pf-c-page__sidebar-body").Body(
			uis...,
		),
	)
}

func pageSession(uis ...func() app.UI) app.UI {
	var el = app.Section().Class("pf-c-page__main-section pf-m-fill pf-m-light")
	return el.Body(jsutil.UIWrap(uis...))
}

func (h *Home) pageHeader(uis ...app.UI) app.UI {
	// page__header
	hd := app.Header().Class("pf-c-page__header", "pf-u-display-flex")
	return hd.Body(jsutil.UIWrap(
		// page__header-brand
		func() app.UI {
			brand := app.Div().Class("pf-c-page__header-brand")
			return brand.Body(
				app.A().
					Href("#").
					Class("pf-c-page__header-brand-link").Body(
					app.Img().
						Class("pf-c-brand").
						Src("/img/favicon.png").
						Alt("Logo"),
				),
			)
		},
		// page__header-tools
		func() app.UI {
			return app.Div().
				Class("pf-c-page__header-tools").Body(
				app.Div().
					Class("pf-c-page__header-tools-group").Body(
					app.Div().
						Class("pf-c-page__header-tools-item").Body(
						app.A().
							Href("https://github.com/pojntfx/html2goapp").
							Target("_blank").
							Class("pf-c-button pf-m-plain").
							Aria("label", "Help").Body(
							app.I().
								Class("pf-icon pf-icon-help").
								Aria("hidden", true),
						),
					),
				),
			)
		}),
	)
}

func (h *Home) serviceUI(uis ...app.UI) app.UI {
	return nil
}

func card() {
	app.Div().
		Class("pf-c-card").
		Body(
			app.Div().
				Class("pf-c-card__title").
				Text("Input"),
			app.Div().
				Class("pf-c-card__body").
				Body())
}

func (h *Home) pageMain(uis ...app.UI) app.UI {
	el := app.Main().
		ID("main").
		Class("pf-c-page__main").
		TabIndex(-1)

	return el.Body(
		pageSession(func() app.UI {
			// list bordered
			var listServiceUI = func() app.UI {
				return app.Ul().Class("pf-c-list pf-m-plain pf-m-bordered").Body(
					app.Range(h.methods).Slice(func(i int) app.UI {
						return app.Li().Class("pf-c-list__item").Body(
							app.A().Class("pf-c-list__item-icon").Body(
								app.I().Class("fas fa-times").
									Aria("hidden", true).OnClick(func(ctx app.Context, e app.Event) {
									fmt.Println("close:", h.methods[i])
								}),
							),

							app.A().Class("pf-c-list__item-text").Body(
								app.Text(h.methods[i])).OnClick(func(ctx app.Context, e app.Event) {
								fmt.Println("srv:", h.methods[i])
								schema, template := h.functionDescribe(h.target, h.methods[i])
								h.input = template
								h.output = schema
								if h.editor == nil {
									h.editor = ace.New("editor")
									h.editor.SetReadOnly(true)
								}
								h.editor.SetValue(h.input)
							}),

							app.A().Class("pf-c-list__item-icon").Body(
								app.I().Class("fas fa-arrow-right").
									Aria("hidden", true).OnClick(func(ctx app.Context, e app.Event) {
									fmt.Println("close:", h.methods[i])

									fmt.Println(h.editor.GetValue())
									h.invokeFunc(h.target, h.methods[i], h.editor.GetValue())
								}),
							),
						)
					}),
				)
			}

			// dropdown
			var dropdown = func() app.UI {
				return app.Div().Class("pf-c-dropdown").Body(
					app.Div().Class("pf-c-input-group").Body(jsutil.UIWrap(
						// Input
						func() app.UI {
							return app.Input().
								Value(h.target).
								Class("pf-c-form-control").
								Name("textarea1").
								ID("textarea1").
								Aria("label", "Textarea with buttons").
								Aria("describedby", "textAreaButton1").OnChange(func(ctx app.Context, e app.Event) {
								h.target = ctx.JSSrc().Get("value").String()
								fmt.Println(h.target)
							})
						},

						// dropdown__toggle
						func() app.UI {
							return app.Div().Class("pf-c-dropdown__toggle pf-m-split-button pf-m-action").Body(jsutil.UIWrap(
								func() app.UI {
									return app.Button().Class("pf-c-dropdown__toggle-button").
										Type("button").
										Aria("label", "Dropdown toggle").
										Body(app.Text("Action")).
										OnClick(func(ctx app.Context, e app.Event) {
											h.services = h.listServices(h.target)
											fmt.Println(h.services)
										})
								},
								func() app.UI {
									return app.Button().
										Class("pf-c-dropdown__toggle-button").
										Type("button").
										Aria("expanded", false).
										ID("dropdown-split-button-action-toggle-button").
										Aria("label", "Dropdown toggle").Body(
										app.I().
											Class("fas fa-caret-down").
											Aria("hidden", true),
									)
								},
							))
						},

						// dropdown__menu
						func() app.UI {
							return app.Ul().
								Class("pf-c-dropdown__menu").
								Hidden(true).Body(jsutil.UIWrap(
								func() app.UI {
									return app.Li().Body(
										app.Button().
											Class("pf-c-dropdown__menu-item").
											Type("button").Body(
											app.Text("Actions"),
										),
									)
								},

								func() app.UI {
									return app.Li().Body(
										app.Button().
											Class("pf-c-dropdown__menu-item").
											Type("button").
											Disabled(true).Body(
											app.Text("Disabled action"),
										),
									)
								},

								func() app.UI {
									return app.Li().Body(
										app.Button().
											Class("pf-c-dropdown__menu-item").
											Type("button").Body(
											app.Text("Other action"),
										),
									)
								},
							))
						},
					)))
			}
			var services = func() app.UI {
				return app.Div().
					Class("pf-c-input-group").Body(
					app.Button().
						Class("pf-c-button pf-m-control").
						Type("button").
						ID("textAreaButton1").Body(
						app.Text("Button"),
					),
					app.Div().
						Class("pf-c-select").Body(
						app.Span().
							ID("select-single-label").
							Hidden(true).
							Body(
								app.Text("Choose one"),
							),
						app.Button().
							Class("pf-c-select__toggle").
							Type("button").
							ID("select-single-toggle").
							Aria("haspopup", true).
							Aria("expanded", "false").
							Aria("labelledby", "select-single-label select-single-toggle").
							Body(
								app.Div().
									Class("pf-c-select__toggle-wrapper").
									Body(
										app.Span().
											Class("pf-c-select__toggle-text").
											Body(
												app.Text(jsutil.IfElse(h.curSrv != "", h.curSrv, "Filter by status")),
											),
									),
								app.Span().
									Class("pf-c-select__toggle-arrow").
									Body(
										app.I().
											Class("fas fa-caret-down").
											Aria("hidden", true),
									)).OnClick(func(ctx app.Context, e app.Event) {
							h.expanded = !h.expanded
						}),
						app.Ul().
							Class("pf-c-select__menu").
							Aria("role", "listbox").
							Aria("labelledby", "select-single-label").
							Hidden(!h.expanded).Body(
							app.Range(h.services).Slice(func(i int) app.UI {
								return app.Li().
									Aria("role", "presentation").Body(
									app.Button().
										Class("pf-c-select__menu-item").
										Aria("role", "option").
										Body(app.Text(h.services[i])).OnClick(func(ctx app.Context, e app.Event) {
										fmt.Println(h.services[i])
										h.curSrv = h.services[i]
										h.methods = h.listFuncs(h.target, h.services[i])
										h.expanded = false

										//	fmt.Println(h.listServices("localhost:8080"))
										//						fmt.Println(h.listActiveSrv())
										//						fmt.Println(h.listFuncs("localhost:8080", "turingvideo.perm.v1.OrgSrv"))
										//						fmt.Println(h.functionDescribe("localhost:8080", "turingvideo.perm.v1.OrgSrv"))
										//						fmt.Println(h.serviceClose("localhost:8080"))
										//						fmt.Println(h.listActiveSrv())
										//						h.invokeFunc("localhost:8080", "turingvideo.perm.v1.OrgSrv.ListOrg", `{
										//  "orgId": "",
										//  "userId": "",
										//  "resType": [
										//    ""
										//  ]
										//}`)
									}),
								)
							}),
						),
					),

					app.Div().
						Class("pf-c-select").Body(
						app.Span().
							ID("select-single-label").
							Hidden(true).
							Body(
								app.Text("Choose one"),
							),
						app.Button().
							Class("pf-c-select__toggle").
							Type("button").
							ID("select-single-toggle").
							Aria("haspopup", true).
							Aria("expanded", "false").
							Aria("labelledby", "select-single-label select-single-toggle").
							Body(
								app.Div().
									Class("pf-c-select__toggle-wrapper").
									Body(
										app.Span().
											Class("pf-c-select__toggle-text").
											Body(
												app.Text(jsutil.IfElse(h.curMth != "", h.curMth, "Filter by status")),
											),
									),
								app.Span().
									Class("pf-c-select__toggle-arrow").
									Body(
										app.I().
											Class("fas fa-caret-down").
											Aria("hidden", true),
									)).OnClick(func(ctx app.Context, e app.Event) {
							h.expanded1 = !h.expanded1
						}),
						app.Ul().
							Class("pf-c-select__menu").
							Aria("role", "listbox").
							Aria("labelledby", "select-single-label").
							Hidden(!h.expanded1).Body(
							app.Range(h.methods).Slice(func(i int) app.UI {
								return app.Li().
									Aria("role", "presentation").Body(
									app.Button().
										Class("pf-c-select__menu-item").
										Aria("role", "option").
										Body(app.Text(h.methods[i])).OnClick(func(ctx app.Context, e app.Event) {
										fmt.Println(h.methods[i])
										h.curMth = h.methods[i]
										h.expanded1 = false
										//h.methods = h.listFuncs(h.target, h.services[i])

										//	fmt.Println(h.listServices("localhost:8080"))
										//						fmt.Println(h.listActiveSrv())
										//						fmt.Println(h.listFuncs("localhost:8080", "turingvideo.perm.v1.OrgSrv"))
										//						fmt.Println(h.functionDescribe("localhost:8080", "turingvideo.perm.v1.OrgSrv"))
										//						fmt.Println(h.serviceClose("localhost:8080"))
										//						fmt.Println(h.listActiveSrv())
										//						h.invokeFunc("localhost:8080", "turingvideo.perm.v1.OrgSrv.ListOrg", `{
										//  "orgId": "",
										//  "userId": "",
										//  "resType": [
										//    ""
										//  ]
										//}`)
									}),
								)
							}),
						),
					),
				)
			}
			var cardForm = func() app.UI {
				return app.Div().
					Class("pf-c-card").
					Body(
						dropdown(),
						services(),
						jsutil.UIWrap(
							func() app.UI {
								return app.Div().Class("pf-c-card__body").Body(
									app.Pre().ID("editor"),
								)
							},

							func() app.UI {
								return app.Table().
									Class("pf-c-table pf-m-grid-md").
									Aria("role", "grid").
									Aria("label", "This is a simple table example").
									ID("table-basic").
									Body(
										app.THead().
											Body(
												app.Tr().
													Aria("role", "row").
													Body(
														app.Th().
															Aria("role", "columnheader").
															Scope("col").
															Body(
																app.Input().
																	Class("pf-c-form-control").
																	Required(true).
																	Type("text"),
															),
														app.Th().
															Aria("role", "columnheader").
															Scope("col").
															Body(
																app.Div().Text("Branches").ContentEditable(true),
															),
														app.Th().
															Aria("role", "columnheader").
															Scope("col").
															Body(
																app.Text("Pull requests"),
															),
													),
											),
										app.TBody().
											Aria("role", "rowgroup").
											Body(
												app.Tr().
													Aria("role", "row").
													Body(
														app.Td().
															Aria("role", "cell").
															DataSet("label", "Repository name").
															Body(
																app.Input(),
															),
														app.Td().
															Aria("role", "cell").
															DataSet("label", "Branches").
															Body(
																app.Input().
																	Class("pf-c-form-control"),
															),
														app.Td().
															Aria("role", "cell").
															DataSet("label", "Pull requests").
															Body(
																app.Input().
																	Class("pf-c-form-control"),
															),
													),
											),
									)
							}),
						app.Div().
							Class("pf-c-card__title").
							Text("Input"),
						app.Div().
							Class("pf-c-card__body").
							Body(
								app.Form().
									Class("pf-c-form").
									OnSubmit(func(ctx app.Context, e app.Event) {
										e.PreventDefault()
									}).Body(jsutil.UIWrap(
									func() app.UI {
										return app.Div().
											Class("pf-c-form__group").Body(
											app.Div().
												Class("pf-c-form__group-label").
												Body(
													app.Label().
														Class("pf-c-form__label").
														For("go-app-pkg-input").
														Body(
															app.Span().
																Class("pf-c-form__label-text").
																Text("go-App Package"),
															app.Span().
																Class("pf-c-form__label-required").
																Aria("hidden", true).
																Text("*"),
														),
												),
											app.Div().
												Class("pf-c-form__group-control").
												Body(
													app.Input().
														Class("pf-c-form-control").
														Required(true).
														OnInput(func(ctx app.Context, e app.Event) {
															if input := ctx.JSSrc().Get("value").String(); input != "" {
															}
														}).
														Type("text").
														ID("go-app-pkg-input"),
												),
										)
									},
									func() app.UI {
										return app.Div().
											Class("pf-c-form__group").
											Body(
												app.Div().
													Class("pf-c-form__group-label").
													Body(
														app.Label().
															Class("pf-c-form__label").
															For("component-pkg-input").
															Body(
																app.Span().
																	Class("pf-c-form__label-text").
																	Text("Target Package"),
																app.Span().
																	Class("pf-c-form__label-required").
																	Aria("hidden", true).
																	Text("*"),
															),
													),
												app.Div().
													Class("pf-c-form__group-control").
													Body(
														app.Input().
															Class("pf-c-form-control").
															Required(true).
															OnInput(func(ctx app.Context, e app.Event) {
																if input := ctx.JSSrc().Get("value").String(); input != "" {
																}
															}).
															Type("text").
															ID("component-pkg-input"),
													),
											)
									},
									func() app.UI {
										return app.Div().
											Class("pf-c-form__group").
											Body(
												app.Div().
													Class("pf-c-form__group-label").
													Body(
														app.Label().
															Class("pf-c-form__label").
															For("component-name-input").
															Body(
																app.Span().
																	Class("pf-c-form__label-text").
																	Text("Component Name"),
																app.Span().
																	Class("pf-c-form__label-required").
																	Aria("hidden", true).
																	Text("*"),
															),
													),
												app.Div().
													Class("pf-c-form__group-control").
													Body(
														app.Input().
															Class("pf-c-form-control").
															Type("text").
															Required(true).
															OnInput(func(ctx app.Context, e app.Event) {
															}).
															Value("c.component").
															ID("component-name-input"),
													),
											)
									},
									func() app.UI {
										return app.Div().
											Class("pf-c-form__group").
											Body(
												app.Div().
													Class("pf-c-form__group-label").
													Body(
														app.Label().
															Class("pf-c-form__label").
															For("html-input").
															Body(
																app.Span().
																	Class("pf-c-form__label-text").
																	Text("Source Code"),
																app.Span().
																	Class("pf-c-form__label-required").
																	Aria("hidden", true).
																	Text("*"),
															),
													),
												app.Div().
													Class("pf-c-form__group-control").
													Body(
														app.Div().
															Class("pf-c-code-editor").
															Body(
																app.Div().
																	Class("pf-c-code-editor__header").
																	Body(jsutil.UIWrap(
																		func() app.UI {
																			return app.Div().
																				Class("pf-c-code-editor__controls").
																				Body(
																					app.Button().
																						Class("pf-c-button pf-m-control").
																						Type("button").
																						Aria("label", "Format").
																						OnClick(func(ctx app.Context, e app.Event) {
																						}).
																						Body(
																							app.I().
																								Class("fas fa-magic").
																								Aria("hidden", true),
																						),
																				)
																		},

																		func() app.UI {
																			return app.Div().
																				Class("pf-c-code-editor__tab").
																				Body(
																					app.Span().
																						Class("pf-c-code-editor__tab-icon").
																						Body(
																							app.I().
																								Class("fas fa-code"),
																						),
																					app.Span().
																						Class("pf-c-code-editor__tab-text").
																						Body(
																							app.Text("HTML"),
																						),
																				)
																		},

																		func() app.UI {
																			return app.Div().
																				Class("pf-c-code-editor__main").
																				Body(
																					app.Pre().ID("editor1"),
																					//app.Textarea().
																					//	ID("html-input").
																					//	Placeholder("Enter HTML input here").
																					//	Required(true).
																					//	//Style("width", "100%").
																					//	//Style("resize", "vertical").
																					//	//Style("border", "0").
																					//	Class("pf-c-form-control").
																					//	Rows(25).
																					//	Text(h.input),
																				)
																		},
																	)),
															),
													),
											)
									},
									func() app.UI {
										return app.Div().
											Class("pf-c-form__group").
											Body(
												app.Div().
													Class("pf-c-form__group-control").
													Body(
														app.Div().
															Class("pf-c-form__actions").
															Body(
																app.Button().
																	Class("pf-c-button pf-m-primary").
																	Type("submit").
																	Text("Convert to Go"),
															),
													),
											)
									},
								)),
							),
					)
			}
			var cardExample = func() app.UI {
				return app.Div().
					Class("pf-c-card pf-m-rounded").
					ID("card-rounded-example").Body(jsutil.UIWrap(
					func() app.UI {
						return app.Div().
							Class("pf-c-card__title").
							Body(
								app.Text("Title"),
							)
					},
					func() app.UI {
						return app.Div().
							Class("pf-c-card__body").
							Body(
								app.Textarea().
									Placeholder("go-app's syntax will be here").
									ReadOnly(true).
									Style("width", "100%").
									Style("resize", "vertical").
									Style("border", "0").
									Class("pf-c-form-control").
									Rows(20).
									Text(h.output),
							)
					},
				))
			}
			var cardEdit = func() app.UI {
				return app.Div().
					Class("pf-c-card").
					Body(
						app.Div().Class("pf-c-card__title").
							Text("Output"),
						app.Div().Class("pf-c-card__body").
							Body(jsutil.UIWrap(
								func() app.UI {
									var editorHeader = func() app.UI {
										return app.Div().
											Class("pf-c-code-editor__header").
											Body(
												app.Div().
													Class("pf-c-code-editor__tab").
													Body(
														app.Span().
															Class("pf-c-code-editor__tab-icon").
															Body(
																app.I().Class("fas fa-code"),
															),
														app.Span().
															Class("pf-c-code-editor__tab-text").Text("Go"),
													),
											)
									}
									var editorBody = func() app.UI {
										return app.Div().
											Class("pf-c-code-editor__main").
											Body(jsutil.UIWrap(
												func() app.UI {
													return app.Textarea().
														Placeholder("go-app's syntax will be here").
														ReadOnly(true).
														//Style("width", "100%").
														//Style("resize", "vertical").
														//Style("border", "0").
														Class("pf-c-form-control").
														Rows(25).
														Text(h.output)
												},
												func() app.UI {
													return app.Div().
														Class("pf-c-code-editor").Body(jsutil.UIWrap(
														func() app.UI {
															return app.Div().
																Class("pf-c-code-editor__header").
																Body(
																	app.Div().
																		Class("pf-c-code-editor__controls").
																		Body(
																			app.Button().
																				Class("pf-c-button pf-m-control").
																				Type("button").
																				Aria("label", "Copy to clipboard").
																				Body(
																					app.I().
																						Class("fas fa-copy").
																						Aria("hidden", true),
																				),
																			app.Button().
																				Class("pf-c-button pf-m-control").
																				Type("button").
																				Aria("label", "Download code").
																				Body(
																					app.I().
																						Class("fas fa-download"),
																				),
																			app.Button().
																				Class("pf-c-button pf-m-control").
																				Type("button").
																				Aria("label", "Upload code").
																				Body(
																					app.I().
																						Class("fas fa-upload"),
																				),
																		),
																	app.Div().
																		Class("pf-c-code-editor__header-main"),
																	app.Div().
																		Class("pf-c-code-editor__tab").
																		Body(
																			app.Span().
																				Class("pf-c-code-editor__tab-icon").
																				Body(
																					app.I().
																						Class("fas fa-code"),
																				),
																			app.Span().
																				Class("pf-c-code-editor__tab-text").
																				Body(
																					app.Text("HTML"),
																				),
																		),
																)
														},
														func() app.UI {
															return app.Div().
																Class("pf-c-code-editor__main").
																Body(
																	app.Code().
																		Class("pf-c-code-editor__code").
																		Body(
																			app.Pre().
																				Class("pf-c-code-editor__code-pre").
																				Body(
																					app.Text(h.output),
																				),
																		),
																)
														},
													))
												},
											))
									}
									return app.Div().Class("pf-c-code-editor pf-m-read-only").
										Body(editorHeader(), editorBody())
								},
							)),
					)
			}

			return grid(
				gridItem(3)(listServiceUI()),
				gridItem(4, 3)(cardForm()),
				gridItem(5, 7)(cardExample()),
				gridItem(6, 3)(cardEdit()),
			)
		}),
	)
}

func (h *Home) Render() app.UI {
	fmt.Println("Render", h.Mounted())

	return page(
		h.pageHeader(),
		//h.pageSidebar(),
		h.pageMain(),
	)
}

func grid(items ...app.UI) app.UI {
	var el = app.Div().Class("pf-l-grid", "pf-m-gutter")
	return el.Body(items...)
}

func gridItem(num ...int) func(items ...app.UI) app.UI {
	var item []string
	item = append(item, "pf-l-grid__item")
	if len(num) > 0 {
		item = append(item, fmt.Sprintf("pf-m-%d-col", num[0]))
	}

	if len(num) > 1 {
		item = append(item, fmt.Sprintf("pf-m-offset-%d-col", num[1]))
	}

	return func(items ...app.UI) app.UI {
		return app.Div().Class(item...).Body(items...)
	}
}
