package app

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pubgo/xerror"
	"honnef.co/go/js/dom/v2"
	"io/ioutil"
	"net/http"

	"github.com/gusaul/grpcox/web/jsutil"
)

var doc = dom.GetWindow().Document()

func (h *Home) OnInit() {
	fmt.Println("OnInit", h.Mounted())
	h.data = make(map[string]*Request)
	h.tables = make(map[string]bool)
	h.picker = (&jsutil.FilePicker{ID: "hiddenFilePicker", Multiple: false}).Accept("image/*")
	h.tableHidden = true
	h.req = new(Request)
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
	fmt.Println("OnMount", h.Mounted())

	fmt.Println("OnMount", h.data)

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

	return app.Div().Body(
		h.container(),
		h.getSource(),
		h.getConnections(),
		h.spinner(),
		h.saveRequestUi(),
	)
}
