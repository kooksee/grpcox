package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fullstorydev/grpchan/httpgrpc"
	"github.com/google/uuid"
	"github.com/gusaul/grpcox/internal/proto/demov1pb"
	"github.com/gusaul/grpcox/web/jsutil"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pubgo/xerror"
	"honnef.co/go/js/dom/v2"
	"io/ioutil"
	fetch "marwan.io/wasm-fetch"
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

func (h *Home) container() app.UI {
	var ss = h.req.Name
	return app.Div().
		Class("container-fluid pt-50").
		Body(
			app.Div().
				Class("row animated fadeIn justify-content-md-center").
				Body(
					app.Div().
						Class("col-3").
						Body(
							app.Div().
								Class("row").
								Body(
									app.Div().
										Class("col").
										Style("padding-left", " 70px;").
										Body(
											app.Div().
												Class("row column-row-left").
												Body(
													app.Div().
														Class("md-form input-group").
														Body(
															app.Input().
																Type("text").
																Class("form-control search").
																ID("search-request").
																OnKeyUp(func(ctx app.Context, e app.Event) {
																	fmt.Println("search-request")
																	app.Script().Text(`
const li = document.querySelectorAll(".request-list")
    li.forEach(function (el) {
        if (el.getAttribute("request-name").toLowerCase().includes(elm.value.toLowerCase())) {
            el.style.display = ""
        } else {
            el.style.display = "none"
        }
    })
`)
																	//	"search(this)"
																}),
															app.Label().
																For("search-request").
																Class("").
																Body(
																	app.Text("Search Request"),
																),
														),
													app.Ul().
														Class("list-group list-group-flush list").
														ID("request-list").Body(
														app.Range(h.data).Map(func(k string) app.UI {
															var node = app.Li().Class("list-group-item").Class("request-list")
															node.Attr("request-name", k)
															node.Body(
																app.A().
																	Title("Delete this request").
																	Class("delete-request").
																	Body(
																		app.I().
																			Class("fa fa-times"),
																	).OnClick(h.removeRequestEvent(h.data[k].ID)),
																app.P().
																	Class("one-long-line request").
																	Body(
																		app.Text(h.data[k].Name),
																	).OnClick(h.updateRequestView(h.data[k].ID)),
															)
															return node
														}),
													),
												),
										),
								),
						),
					app.Div().
						Class("col-7").
						Body(
							app.Div().
								Class("md-form input-group").
								Body(
									app.Input().
										Type("text").
										Class("form-control").Attr("value", ss).
										ID("server-target"),
									app.Label().
										For("server-target").
										Body(
											app.Text("gRPC Server Target"),
										),
									app.Div().
										Class("input-group-append").
										Body(
											app.Button().
												ID("get-services").OnClick(func(ctx app.Context, e app.Event) {
												fmt.Println("click get-services")
												fmt.Println(ctx, e)
											}).
												Class("btn btn-mdb-color waves-effect m-0").
												Type("button").
												Body(
													app.I().
														Class("fa fa-plug"),
												),
											app.Div().
												Class("dropdown").
												Body(
													app.Button().
														Type("button").
														Class("btn btn-mdb-color waves-effect m-0 dropdown-toggle save-button-dropdown").
														DataSet("toggle", "dropdown").
														Aria("haspopup", true).
														Aria("expanded", true),
													app.Div().
														Class("dropdown-menu dropdown-menu-right").
														Aria("labelledby", "btnGroupDrop1").
														Body(
															app.A().
																Class("dropdown-item").
																ID("show-modal-save-request").
																Body(
																	app.Text("Save"),
																).OnClick(h.showModalSaveRequest()),
															app.A().
																Class("dropdown-item").
																ID("show-modal-save-as-request").
																Body(
																	app.Text("Save As"),
																),
														),
												),
										),
								),
							app.Div().
								Class("custom-control custom-checkbox").
								Body(
									app.Input().
										Type("checkbox").
										Class("custom-control-input").
										ID("use-tls"),
									app.Label().
										Class("custom-control-label").
										For("use-tls").
										Body(
											app.Text("Use TLS"),
										),
								),
							app.Div().
								Class("custom-control custom-checkbox").
								Body(
									app.Input().
										Type("checkbox").
										Class("custom-control-input").
										ID("restart-conn"),
									app.Label().
										Class("custom-control-label").
										For("restart-conn").
										Body(
											app.Text("Restart Connection"),
										),
								),
							app.Div().
								Class("input-group").
								Body(
									app.Div().
										Class("custom-control custom-checkbox").
										Body(
											app.Input().
												Type("checkbox").
												Class("custom-control-input").
												ID("local-proto"),
											app.Label().
												Class("custom-control-label").
												For("local-proto").
												Body(
													app.Text("Use local proto"),
												),
										),
								),
							app.Div().
								Class("input-group").
								ID("proto-input").
								Style("display", " none").
								Body(
									app.Div().
										Class("proto-top-collection").
										Body(
											app.Input().
												Class("proto-uploader").
												Type("file").
												ID("proto-file").
												Multiple(true),
											app.Label().
												For("proto-file").
												Body(
													app.I().
														Class("fa fa-plus-circle"),
													app.Text("proto files"),
												),
											app.Span().
												ID("proto-collection-toggle").
												Class("proto-toggle").
												Body(
													app.Text("Hide Proto Collection"),
												),
										),
									app.Div().
										Class("proto-collection"),
								),
							app.Div().
								Class("input-group").
								Body(
									app.Div().
										Class("custom-control custom-checkbox").
										Body(
											app.Input().
												Type("checkbox").
												Class("custom-control-input").
												ID("ctx-metadata-switch").
												OnChange(func(ctx app.Context, e app.Event) {
													h.tableHidden = !e.Get("target").Get("checked").Bool()
													h.toggleDisplayCtxMetadataTable(e.Get("target").Get("checked").Bool())
												}),
											app.Label().
												Class("custom-control-label").
												For("ctx-metadata-switch").
												Body(
													app.Text("Use request metadata"),
												),
										),
								),
							app.Div().
								Class("input-group").
								ID("ctx-metadata-input").
								Style("display", " block").Hidden(h.tableHidden).
								Body(
									app.Br(),
									h.metadataTable(),
								),
							app.Div().
								Class("other-elem").
								ID("choose-service").
								Style("display", " none").
								Body(
									app.Div().
										Class("input-group").
										Body(
											app.Div().
												Class("input-group-prepend").
												Body(
													app.Span().
														Class("input-group-text btn-dark w-120").
														Attr("for", "select-service").
														Body(
															app.I().
																Class("fa fa-television"),
															app.Text("Services"),
														),
												),
											app.Select().
												Class("browser-default custom-select").
												ID("select-service").Body(app.Option().Text("Choose Service")),
										),
								),
							app.Div().
								Class("other-elem").
								ID("choose-function").
								Style("display", " none").
								Body(
									app.Div().
										Class("input-group").
										Body(
											app.Div().
												Class("input-group-prepend").
												Body(
													app.Span().
														Class("input-group-text btn-dark w-120").
														Attr("for", "select-function").
														Body(
															app.I().
																Class("fa fa-rocket"),
															app.Text("Methods"),
														),
												),
											app.Select().
												Class("browser-default custom-select").
												ID("select-function"),
										),
								),
							app.Div().
								Class("row other-elem").
								ID("body-request").
								Style("display", " none").
								Body(
									app.Div().
										Class("col-md-7").
										Body(
											app.Div().
												Class("card").
												Body(
													app.Div().
														Class("card-body schema-body").
														Body(
															app.Pre().
																ID("editor"),
														),
												),
											app.Button().
												Class("btn btn-primary waves-effect mt-10").
												ID("invoke-func").
												Type("button").
												Body(
													app.I().
														Class("fa fa-play"),
													app.Text("Submit"),
												),
										),
									app.Div().
										Class("col-md-5").
										Body(
											app.Div().
												Class("card").
												Body(
													app.Div().
														Class("card-body schema-body").
														Body(
															app.H4().
																Class("card-title").
																Body(
																	app.A().
																		Body(
																			app.Text("Schema Input"),
																		),
																),
															app.Pre().
																Class("prettyprint custom-pretty").
																ID("schema-proto"),
														),
												),
										),
								),
							app.Div().
								Class("row other-elem").
								ID("response").
								Style("display", " none").
								Body(
									app.Div().
										Class("col").
										Body(
											app.Div().
												Class("card").
												Body(
													app.Div().
														Class("card-body").
														Body(
															app.Small().
																Class("pull-right").
																ID("timer-resp").
																Body(
																	app.Text("Time :"), app.Span(),
																),
															app.H4().
																Class("card-title").
																Body(
																	app.A().
																		Body(
																			app.Text("Response:"),
																		),
																),
															app.P().
																Class("card-text").
																Body(),
															app.Pre().
																Class("prettyprint custom-pretty").
																ID("json-response"),
															app.P(),
														),
												),
										),
								),
						),
				),
		)
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
								Class("modal-header").
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
	h.target = "localhost:8080"
	fmt.Println("OnMount", h.Mounted())

	fmt.Println("OnMount", h.data)

	rsp, err := fetch.Fetch(fmt.Sprintf("/js/ace.js"), &fetch.Opts{Method: fetch.MethodGet})
	xerror.Panic(err)
	fmt.Println(string(rsp.Body))
	fmt.Println(jsutil.Eval(string(rsp.Body)))

	if h.editor == nil {
		h.editor = h.ace("editor")
	}

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

	ctx.Async(func() {
		defer xerror.RecoverAndExit()
		var ret, err = h.cc.ServerStream(context.Background(), &demov1pb.Message{
			Hello: "world",
		})
		if err != nil {
			return
		}

		for {
			mm, err := ret.Recv()
			xerror.Panic(err)
			if err != nil {
				break
			} else {
				fmt.Println("ServerStream value:", mm)
			}
		}
	})
}

func (h *Home) OnNav(ctx app.Context) {
	fmt.Println("OnNav", h.Mounted())
	p := ctx.Page()
	p.SetTitle("Image to Factorio Blueprint")
	p.SetAuthor("mlctrez")
	p.SetKeywords("factorio, blueprint, image, convert")
	p.SetDescription("A progressive web application for converting images to factorio tile blueprints.")
}

func (h *Home) Render() app.UI {
	if h.picker == nil {
		h.picker = (&jsutil.FilePicker{ID: "hiddenFilePicker", Multiple: false}).Accept("image/*")
	}

	fmt.Println("Render", h.Mounted())

	page := app.Div().Class("pf-c-page")
	return jsutil.Compo(page).Body(
		// page__header
		func() app.UI {
			hd := app.Header().Class("pf-c-page__header", "pf-u-display-flex")
			return jsutil.Compo(hd).Body(
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
						))
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
				},
			)
		},
		// page__main
		func() app.UI {
			el := app.Main().
				ID("main").
				Class("pf-c-page__main").
				TabIndex(-1)

			return jsutil.Compo(el).Body(
				// dropdown
				func() app.UI {
					return app.Div().
						Class("pf-c-dropdown").Body(
						app.Div().
							Class("pf-c-input-group").Body(jsutil.UI{
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
								return app.Div().
									Class("pf-c-dropdown__toggle pf-m-split-button pf-m-action").Body(
									app.Button().
										Class("pf-c-dropdown__toggle-button").
										Type("button").
										Aria("label", "Dropdown toggle").
										Body(
											app.Text("Action"),
										).OnClick(func(ctx app.Context, e app.Event) {
										h.services = h.listServices(h.target)
										fmt.Println(h.services)
										schema, template := h.functionDescribe("localhost:8080", "turingvideo.perm.v1.OrgSrv")
										h.input = template
										h.output = schema
										if h.editor == nil {
											h.editor = h.ace("editor")
										}
										h.editor.Call("setValue", h.input)
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
									app.Button().
										Class("pf-c-dropdown__toggle-button").
										Type("button").
										Aria("expanded", h.expanded).
										ID("dropdown-split-button-action-toggle-button").
										Aria("label", "Dropdown toggle").
										Body(
											app.I().
												Class("fas fa-caret-down").
												Aria("hidden", true),
										).OnClick(func(ctx app.Context, e app.Event) {
										h.expanded = !h.expanded
										fmt.Println(h.expanded)
										fmt.Println("hidden", jsutil.Hidden(!h.expanded))
									}),
								)
							},
							// dropdown__menu
							func() app.UI {
								return app.Ul().
									Class("pf-c-dropdown__menu").
									Hidden(!h.expanded).Body(
									app.Li().Body(
										app.Button().
											Class("pf-c-dropdown__menu-item").
											Type("button").Body(
											app.Text("Actions"),
										),
									),
									app.Li().Body(
										app.Button().
											Class("pf-c-dropdown__menu-item").
											Type("button").
											Disabled(true).Body(
											app.Text("Disabled action"),
										),
									),
									app.Li().Body(
										app.Button().
											Class("pf-c-dropdown__menu-item").
											Type("button").Body(
											app.Text("Other action"),
										),
									),
								)
							},
						}.Render()))
				},
				// list bordered
				func() app.UI {
					return app.Ul().
						Class("pf-c-list pf-m-plain pf-m-bordered").Body(
						app.Range(h.services).Slice(func(i int) app.UI {
							return app.Li().Body(
								app.Span().
									Class("pf-c-list__item-icon").Body(
									app.I().
										Class("fas fa-times").
										Aria("hidden", true).OnClick(func(srv string) func(ctx app.Context, e app.Event) {
										return func(ctx app.Context, e app.Event) {
											fmt.Println(srv)
										}
									}(h.services[i])),
								),

								app.Span().Body(app.Text(h.services[i])),
							)
						}),
					)
				},
				func() app.UI {
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
				},
				func() app.UI {
					el := app.Section().Class("pf-c-page__main-section pf-m-fill")
					return el.Body(jsutil.UI{
						func() app.UI {
							return app.Div().
								Class("pf-l-grid pf-m-gutter").Body(jsutil.UI{
								func() app.UI {
									return app.Div().
										Class("pf-l-grid__item", "pf-m-4-col pf-m-offset-3-col").
										Body(
											app.Div().
												Class("pf-c-card").
												Body(
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
																}).
																Body(
																	app.Div().
																		Class("pf-c-form__group").
																		Body(
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
																		),
																	app.Div().
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
																		),
																	app.Div().
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
																		),
																	app.Div().
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
																								Body(
																									app.Div().
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
																										),
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
																								),
																							app.Div().
																								Class("pf-c-code-editor__main").
																								Body(
																									app.Pre().ID("editor"),
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
																								),
																						),
																				),
																		),
																	app.Div().
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
																		),
																),
														),
												),
										)
								},
								func() app.UI {
									return app.Div().
										Class("pf-l-grid__item", "pf-m-2-col pf-m-offset-7-col").
										Body(
											app.Div().
												Class("pf-c-card pf-m-rounded").
												ID("card-rounded-example").
												Body(
													app.Div().
														Class("pf-c-card__title").
														Body(
															app.Text("Title"),
														),
													app.Div().
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
														),
												),
										)
								},
							}.Render())
						},
						func() app.UI {
							return app.Div().
								Class("pf-l-grid pf-m-gutter").
								Body(
									app.Div().
										Class("pf-l-grid__item pf-m-6-col pf-m-offset-3-col").
										Body(
											app.Div().
												Class("pf-c-card").
												Body(
													app.Div().
														Class("pf-c-card__title").
														Text("Output"),
													app.Div().
														Class("pf-c-card__body").
														Body(
															app.Div().
																Class("pf-c-code-editor pf-m-read-only").
																Body(
																	app.Div().
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
																		),
																	app.Div().
																		Class("pf-c-code-editor__main").
																		Body(
																			app.Textarea().
																				Placeholder("go-app's syntax will be here").
																				ReadOnly(true).
																				//Style("width", "100%").
																				//Style("resize", "vertical").
																				//Style("border", "0").
																				Class("pf-c-form-control").
																				Rows(25).
																				Text(h.output),

																			jsutil.UI{
																				func() app.UI {
																					return app.Div().
																						Class("pf-c-code-editor").
																						Body(
																							app.Div().
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
																								),
																							app.Div().
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
																								),
																						)
																				},
																			}.Render()),
																),
														),
												),
										),
								)
						},
					}.Render())
				},
			)
		},
	)
}
