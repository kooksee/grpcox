package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fullstorydev/grpchan/httpgrpc"
	"github.com/google/uuid"
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

func (h *Home) spinner() app.UI {
	return app.Div().
		Class("spinner").
		Style("display", " none").
		Body(
			app.Div().
				Class("rect1"),
			app.Div().
				Class("rect2"),
			app.Div().
				Class("rect3"),
			app.Div().
				Class("rect4"),
			app.Div().
				Class("rect5"),
		)
}

func (h *Home) saveRequestUi() app.UI {
	return app.Div().
		Class("modal fade").
		ID("saveRequest").
		TabIndex(-1).
		Aria("role", "dialog").
		Aria("labelledby", "exampleModalLabel").
		Aria("hidden", true).
		Body(
			app.Div().
				Class("modal-dialog").
				Aria("role", "document").
				Body(
					app.Div().
						Class("modal-content").
						Body(
							app.Div().
								Class("modal-pageHeader").
								Body(
									app.H5().
										Class("modal-title").
										ID("exampleModalLabel").
										Body(
											app.Text("Input the name for the request"),
										),
								),
							app.Div().
								Class("modal-body").
								Body(
									app.Form().
										Body(
											app.Div().
												Class("form-group row").
												Body(
													app.Label().
														For("input-request-name").
														Class("col-sm-2 col-form-label").
														Body(
															app.Text("Name"),
														),
													app.Div().
														Class("col-sm-10").
														Body(
															app.Input().
																Type("text").
																Class("form-control").
																ID("input-request-name"),
														),
												),
										),
								),
							app.Div().
								Class("modal-footer").
								Body(
									app.Button().
										Type("button").
										Class("btn btn-secondary").
										DataSet("dismiss", "modal").
										Body(
											app.Text("Cancel"),
										),
									app.Button().
										ID("save-request").
										Type("button").
										Class("btn btn-primary").
										Body(
											app.Text("Save"),
										).OnClick(h.saveRequest()),
								),
						),
				),
		)
}

func (h *Home) getConnections() app.UI {
	return app.Div().
		Class("connections").
		Body(
			app.Div().
				Class("title").
				Body(app.Raw(`<svg class="dots" expanded="true" height="100px" width="100px">
                <circle cx="50%" cy="50%" r="7px"></circle>
                <circle class="pulse" cx="50%" cy="50%" r="10px"></circle>
            </svg>`),
					app.Span(),
					app.Text("Active Connection(s)")),
			app.Div().
				ID("conn-list-template").
				Style("display", "none").
				Body(
					app.Li().
						Body(
							app.I().
								Class("fa fa-close").
								DataSet("toggle", "tooltip").
								Title("false"),
							app.Span().
								Class("ip"),
						),
				),
			app.Ul().
				Class("nav"),
		)
}

func (h *Home) toggleDisplayCtxMetadataTable(show bool) {
	var style = "display: none"
	if show {
		style = "display: block"
	}

	var e = dom.GetWindow().Document().GetElementByID("ctx-metadata-input")
	e.RemoveAttribute("style")
	e.SetAttribute("style", style)
}

func (h *Home) newTr(id string) app.UI {
	return app.Tr().
		Body(
			app.Td().
				Body(
					app.Span().
						Class("table-remove").
						Body(
							app.Button().ID(id).OnClick(func(ctx app.Context, e app.Event) {
								var ee = jsutil.Event(e)
								delete(h.tables, ee.CurrentTarget().ID())
								fmt.Println(ee.CurrentTarget().ID())
							}).
								Type("button").
								Class("btn btn-danger btn-rounded btn-sm my-0").
								Body(
									app.I().
										Class("fa fa-times"),
								),
						),
				),
			app.Td().
				Class("ctx-metadata-input-field pt-3-half").
				ContentEditable(true),
			app.Td().
				Class("ctx-metadata-input-field pt-3-half").
				ContentEditable(true),
		)
}

func (h *Home) metadataTable() app.UI {
	fmt.Println("tables=>", h.tables)

	return app.Div().
		ID("ctx-metadata-table").
		Class("table-editable").
		Body(
			app.Table().
				Class("table table-bordered").
				Body(
					app.THead().
						Body(
							app.Tr().
								Body(
									app.Th().
										Class("text-start").
										Style("width", " 10%"),
									app.Th().
										Class("text-start").
										Style("width", " 20%").
										Body(
											app.Text("Key"),
										),
									app.Th().
										Class("text-start").
										Style("width", " 70%").
										Body(
											app.Text("Value"),
										),
								),
						),
					app.TBody().
						Body(app.Range(h.tables).Map(func(s string) app.UI {
							return h.newTr(s)
						})),
				),
			app.Div().
				Class("input-group-append").
				Body(
					app.Span().
						Class("table-add").
						Body(
							app.Button().OnClick(func(ctx app.Context, e app.Event) {
								var id = uuid.New().String()
								fmt.Println("uuid", id)
								h.tables[id] = true
							}).
								Type("button").
								Class("btn btn-success btn-rounded btn-sm my-0").
								Body(
									app.I().
										Class("fa fa-plus"),
								),
						),
				),
		)
}

func (h *Home) getSource() app.UI {
	return app.
		A().
		Class("github-corner").
		Target("_blank").
		Href("https://github.com/gusaul/grpcox").
		Aria("label", "View source on GitHub").
		Body(app.Raw(
			`<svg width="80" height="80" viewBox="0 0 250 250" style="fill:#151513; color:#fff; position: absolute; top: 0; border: 0; right: 0;"
      aria-hidden="true">
      <path d="M0,0 L115,115 L130,115 L142,142 L250,250 L250,0 Z"></path>
      <path d="M128.3,109.0 C113.8,99.7 119.0,89.6 119.0,89.6 C122.0,82.7 120.5,78.6 120.5,78.6 C119.2,72.0 123.4,76.3 123.4,76.3 C127.3,80.9 125.5,87.3 125.5,87.3 C122.9,97.6 130.6,101.9 134.4,103.2"
        fill="currentColor" style="transform-origin: 130px 106px;" class="octo-arm"></path>
      <path d="M115.0,115.0 C114.9,115.1 118.7,116.5 119.8,115.4 L133.7,101.6 C136.9,99.2 139.9,98.4 142.2,98.6 C133.8,88.0 127.5,74.4 143.8,58.0 C148.5,53.4 154.0,51.2 159.7,51.0 C160.3,49.4 163.2,43.6 171.4,40.1 C171.4,40.1 176.1,42.5 178.8,56.2 C183.1,58.6 187.2,61.8 190.9,65.4 C194.5,69.0 197.7,73.2 200.1,77.6 C213.8,80.2 216.3,84.9 216.3,84.9 C212.7,93.1 206.9,96.0 205.4,96.6 C205.1,102.4 203.0,107.8 198.3,112.5 C181.9,128.9 168.3,122.5 157.7,114.1 C157.9,116.9 156.7,120.9 152.7,124.9 L141.0,136.5 C139.8,137.7 141.6,141.9 141.8,141.8 Z"
        fill="currentColor" class="octo-body"></path>
    </svg>`))
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
	var el = app.Section().Class("pf-c-page__main-section pf-m-fill")
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
					app.Range(h.services).Slice(func(i int) app.UI {
						return app.Li().Class("pf-c-list__item").Body(
							app.A().Class("pf-c-list__item-icon").Body(
								app.I().Class("fas fa-times").
									Aria("hidden", true).OnClick(func(ctx app.Context, e app.Event) {
									fmt.Println("close:", h.services[i])
								}),
							),

							app.A().Class("pf-c-list__item-text").Body(
								app.Text(h.services[i])).OnClick(func(ctx app.Context, e app.Event) {
								fmt.Println("srv:", h.services[i])
								schema, template := h.functionDescribe(h.target, h.services[i])
								h.input = template
								h.output = schema
								if h.editor == nil {
									h.editor = ace.New("editor")
									h.editor.SetReadOnly(true)
								}
								h.editor.SetValue(h.input)
							}),
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
										Body(
											app.Text("Action"),
										).OnClick(func(ctx app.Context, e app.Event) {
										h.services = h.listServices(h.target)
										fmt.Println(h.services)

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

									})

								},
								func() app.UI {
									return app.Button().
										Class("pf-c-dropdown__toggle-button").
										Type("button").
										Aria("expanded", h.expanded).
										ID("dropdown-split-button-action-toggle-button").
										Aria("label", "Dropdown toggle").Body(
										app.I().
											Class("fas fa-caret-down").
											Aria("hidden", true),
									).OnClick(func(ctx app.Context, e app.Event) {
										h.expanded = !h.expanded
										fmt.Println(h.expanded)
										fmt.Println("hidden", jsutil.Hidden(!h.expanded))
									})
								},
							))
						},

						// dropdown__menu
						func() app.UI {
							return app.Ul().
								Class("pf-c-dropdown__menu").
								Hidden(!h.expanded).Body(jsutil.UIWrap(
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
												app.Text("Filter by status"),
											),
									),
								app.Span().
									Class("pf-c-select__toggle-arrow").
									Body(
										app.I().
											Class("fas fa-caret-down").
											Aria("hidden", true),
									),
							).OnClick(func(ctx app.Context, e app.Event) {
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
										Body(
											app.Text(h.services[i]),
										),
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
