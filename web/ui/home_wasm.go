package app

import (
	"context"
	"fmt"
	"github.com/fullstorydev/grpchan/httpgrpc"
	"github.com/gusaul/grpcox/internal/proto/demov1pb"
	"github.com/gusaul/grpcox/web/ace"
	"github.com/gusaul/grpcox/web/jsutil"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"honnef.co/go/js/dom/v2"
	"net/http"
	"net/url"
	"strings"
)

var doc = dom.GetWindow().Document()

func (h *Home) OnInit() {
	fmt.Println("OnInit", h.Mounted())
	h.data = make(map[string]*Request)
	h.tables = make(map[string]bool)
	h.picker = (&jsutil.FilePicker{ID: "hiddenFilePicker", Multiple: false}).Accept("image/*")
	h.tableHidden = true
	h.curReq = new(Request)

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

	ctx.Async(func() {
		h.allRequests = h.getAllRequestKey()
		for _, r := range h.allRequests {
			h.data[r.ID] = r
		}
		ctx.Dispatch(func(c app.Context) {})
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
	return app.Div().
		Class("pf-c-page__sidebar  pf-m-light pf-m-expanded").
		Body(
			app.Div().
				Class("pf-c-page__sidebar-body").
				Body(
					app.Nav().
						Class("pf-c-nav pf-m-light").
						ID("nav-expandable-example-expandable-nav").
						Aria("label", "Global").
						Body(
							app.Ul().
								Class("pf-c-nav__list").
								Body(
									app.Li().
										Class("pf-c-nav__item pf-m-expandable pf-m-expanded pf-m-current").
										Body(
											app.Button().
												Class("pf-c-nav__link").
												ID("nav-expandable-example-expandable-nav-link1").
												Aria("expanded", true).
												Body(
													app.Text("System panel"), app.Span().
														Class("pf-c-nav__toggle").
														Body(
															app.Span().
																Class("pf-c-nav__toggle-icon").
																Body(
																	app.I().
																		Class("fas fa-angle-right").
																		Aria("hidden", true),
																),
														),
												),
											app.Section().
												Class("pf-c-nav__subnav").
												Aria("labelledby", "nav-expandable-example-expandable-nav-link1").
												Body(
													app.Ul().
														Class("pf-c-nav__list").
														Body(
															app.Li().
																Class("pf-c-nav__item").
																Body(

																	app.A().
																		Href("#").
																		Class("pf-c-nav__link").
																		Body(
																			app.I().Class("fas fa-times").
																				Aria("hidden", true).OnClick(func(ctx app.Context, e app.Event) {
																				fmt.Println("close:")
																			}),

																			app.Text("Resource usage"),
																		),
																),
															app.Li().
																Class("pf-c-nav__item").
																Body(
																	app.A().
																		Href("#").
																		Class("pf-c-nav__link pf-m-current").
																		Aria("current", "page").
																		Body(
																			app.Text("Resource usage"),
																		),
																),
															app.Li().
																Class("pf-c-nav__item").
																Body(
																	app.A().
																		Href("#").
																		Class("pf-c-nav__link").
																		Body(
																			app.Text("Hypervisors"),
																		),
																),
															app.Li().
																Class("pf-c-divider").
																Aria("role", "separator"),
															app.Li().
																Class("pf-c-nav__item").
																Body(
																	app.A().
																		Href("#").
																		Class("pf-c-nav__link").
																		Body(
																			app.Text("Instances"),
																		),
																),
															app.Li().
																Class("pf-c-nav__item").
																Body(
																	app.A().
																		Href("#").
																		Class("pf-c-nav__link").
																		Body(
																			app.Text("Volumes"),
																		),
																),
															app.Li().
																Class("pf-c-nav__item").
																Body(
																	app.A().
																		Href("#").
																		Class("pf-c-nav__link").
																		Body(
																			app.Text("Networks"),
																		),
																),
														),
												),
										),
									app.Li().
										Class("pf-c-nav__item pf-m-expandable").
										Body(
											app.Button().
												Class("pf-c-nav__link").
												ID("nav-expandable-example-expandable-nav-link2").
												Aria("expanded", "false").
												Body(
													app.Text("Policy"), app.Span().
														Class("pf-c-nav__toggle").
														Body(
															app.Span().
																Class("pf-c-nav__toggle-icon").
																Body(
																	app.I().
																		Class("fas fa-angle-right").
																		Aria("hidden", true),
																),
														),
												),
											app.Section().
												Class("pf-c-nav__subnav").
												Aria("labelledby", "nav-expandable-example-expandable-nav-link2").
												Hidden(true).
												Body(
													app.Ul().
														Class("pf-c-nav__list").
														Body(
															app.Li().
																Class("pf-c-nav__item").
																Body(
																	app.A().
																		Href("#").
																		Class("pf-c-nav__link").
																		Body(
																			app.Text("Subnav link 1"),
																		),
																),
															app.Li().
																Class("pf-c-nav__item").
																Body(
																	app.A().
																		Href("#").
																		Class("pf-c-nav__link").
																		Body(
																			app.Text("Subnav link 2"),
																		),
																),
														),
												),
										),
									app.Li().
										Class("pf-c-nav__item pf-m-expandable").
										Body(
											app.Button().
												Class("pf-c-nav__link").
												ID("nav-expandable-example-expandable-nav-link3").
												Aria("expanded", "false").
												Body(
													app.Text("Authentication"), app.Span().
														Class("pf-c-nav__toggle").
														Body(
															app.Span().
																Class("pf-c-nav__toggle-icon").
																Body(
																	app.I().
																		Class("fas fa-angle-right").
																		Aria("hidden", true),
																),
														),
												),
											app.Section().
												Class("pf-c-nav__subnav").
												Aria("labelledby", "nav-expandable-example-expandable-nav-link3").
												Hidden(true).
												Body(
													app.Ul().
														Class("pf-c-nav__list").
														Body(
															app.Li().
																Class("pf-c-nav__item").
																Body(
																	app.A().
																		Href("#").
																		Class("pf-c-nav__link").
																		Body(
																			app.Text("Subnav link 1"),
																		),
																),
															app.Li().
																Class("pf-c-nav__item").
																Body(
																	app.A().
																		Href("#").
																		Class("pf-c-nav__link").
																		Body(
																			app.Text("Subnav link 2"),
																		),
																),
														),
												),
										),
								),
						),
				),
		)
}

func SkipContent() app.UI {
	return app.A().
		Class("pf-c-skip-to-content pf-c-button pf-m-primary").
		Href("#main-content-drawer-jump-links").Body(
		app.Text("Skip to content"),
	)
}

func pageSession(uis ...func() app.UI) app.UI {
	var el = app.Section().Class("pf-c-page__main-section pf-m-fill pf-m-light")
	return el.Body(jsutil.UIWrap(uis...))
}

func pageHeader(uis ...app.UI) app.UI {
	// page__header
	hd := app.Header().Class("pf-c-page__header")
	return hd.Body(jsutil.UIWrap(
		func() app.UI {
			brand := app.Div().Class("pf-c-page__header-brand")
			return brand.Body(
				app.Script().Src("https://ace.c9.io/build/src/ace.js"),
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
		// tools
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

func drawer() app.UI {
	return app.Raw(`<div class="pf-c-drawer">
      <div class="pf-c-drawer__main">
        <div class="pf-c-drawer__content" id="drawer-jump-links-drawer-scrollable-container">
          <div class="pf-c-drawer__body pf-m-padding">
            <div class="pf-c-sidebar">
              <div class="pf-c-sidebar__main">
                <div class="pf-c-sidebar__panel pf-m-sticky pf-m-gutter">
                  <nav class="pf-c-jump-links pf-m-vertical pf-m-non-expandable-on-md pf-m-expandable">
                    <div class="pf-c-jump-links__label">Jump to section</div>
                    <ul class="pf-c-jump-links__list">
                      <li class="pf-c-jump-links__item pf-m-current">
                        <a class="pf-c-jump-links__link" href="#drawer-jump-links-jump-links-first">
                          <span class="pf-c-jump-links__link-text">First section</span>
                        </a>
                      </li>
                      <li class="pf-c-jump-links__item">
                        <a class="pf-c-jump-links__link" href="#drawer-jump-links-jump-links-second">
                          <span class="pf-c-jump-links__link-text">Second section</span>
                        </a>
                      </li>
                      <li class="pf-c-jump-links__item">
                        <a class="pf-c-jump-links__link" href="#drawer-jump-links-jump-links-third">
                          <span class="pf-c-jump-links__link-text">Third section</span>
                        </a>
                      </li>
                      <li class="pf-c-jump-links__item">
                        <a class="pf-c-jump-links__link" href="#drawer-jump-links-jump-links-fourth">
                          <span class="pf-c-jump-links__link-text">Fourth section</span>
                        </a>
                      </li>
                      <li class="pf-c-jump-links__item">
                        <a class="pf-c-jump-links__link" href="#drawer-jump-links-jump-links-fifth">
                          <span class="pf-c-jump-links__link-text">Fifth section</span>
                        </a>
                      </li>
                    </ul>
                  </nav>
                </div>
                <div class="pf-c-sidebar__content">
                  <div class="pf-c-content">
                    <p>
                      <b>Note:</b> Jump links require javascript to look and function properly, so while clicking on the jump links in these demos may take you to anchors on the page, the full functionality isn't implemented. Refer to the react demos to see the intended design and functionality.
                    </p>

                    <h2 id="drawer-jump-links-jump-links-first">First section</h2>
                    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas sed dui ullamcorper, dignissim purus eu, mattis leo. Curabitur eleifend turpis ipsum, aliquam pretium risus efficitur vel. Etiam eget enim vitae quam facilisis pharetra at eget diam. Suspendisse ut vulputate magna. Nunc viverra posuere orci sit amet pulvinar. Quisque dui justo, egestas ac accumsan suscipit, tristique eu risus. In aliquet libero ante, ac pulvinar lectus pretium in. Ut enim tellus, vulputate et lorem et, hendrerit rutrum diam. Cras pharetra dapibus elit vitae ullamcorper. Nulla facilisis euismod diam, at sodales sem laoreet hendrerit. Curabitur porttitor ex nulla. Proin ligula leo, egestas ac nibh a, pellentesque mollis augue. Donec nec augue vehicula eros pulvinar vehicula eget rutrum neque. Duis sit amet interdum lorem. Vivamus porttitor nec quam a accumsan. Nam pretium vitae leo vitae rhoncus.</p>

                    <h2 id="drawer-jump-links-jump-links-second">Second section</h2>
                    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas sed dui ullamcorper, dignissim purus eu, mattis leo. Curabitur eleifend turpis ipsum, aliquam pretium risus efficitur vel. Etiam eget enim vitae quam facilisis pharetra at eget diam. Suspendisse ut vulputate magna. Nunc viverra posuere orci sit amet pulvinar. Quisque dui justo, egestas ac accumsan suscipit, tristique eu risus. In aliquet libero ante, ac pulvinar lectus pretium in. Ut enim tellus, vulputate et lorem et, hendrerit rutrum diam. Cras pharetra dapibus elit vitae ullamcorper. Nulla facilisis euismod diam, at sodales sem laoreet hendrerit. Curabitur porttitor ex nulla. Proin ligula leo, egestas ac nibh a, pellentesque mollis augue. Donec nec augue vehicula eros pulvinar vehicula eget rutrum neque. Duis sit amet interdum lorem. Vivamus porttitor nec quam a accumsan. Nam pretium vitae leo vitae rhoncus.</p>

                    <h2 id="drawer-jump-links-jump-links-third">Third section</h2>
                    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas sed dui ullamcorper, dignissim purus eu, mattis leo. Curabitur eleifend turpis ipsum, aliquam pretium risus efficitur vel. Etiam eget enim vitae quam facilisis pharetra at eget diam. Suspendisse ut vulputate magna. Nunc viverra posuere orci sit amet pulvinar. Quisque dui justo, egestas ac accumsan suscipit, tristique eu risus. In aliquet libero ante, ac pulvinar lectus pretium in. Ut enim tellus, vulputate et lorem et, hendrerit rutrum diam. Cras pharetra dapibus elit vitae ullamcorper. Nulla facilisis euismod diam, at sodales sem laoreet hendrerit. Curabitur porttitor ex nulla. Proin ligula leo, egestas ac nibh a, pellentesque mollis augue. Donec nec augue vehicula eros pulvinar vehicula eget rutrum neque. Duis sit amet interdum lorem. Vivamus porttitor nec quam a accumsan. Nam pretium vitae leo vitae rhoncus.</p>

                    <h2 id="drawer-jump-links-jump-links-fourth">Fourth section</h2>
                    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas sed dui ullamcorper, dignissim purus eu, mattis leo. Curabitur eleifend turpis ipsum, aliquam pretium risus efficitur vel. Etiam eget enim vitae quam facilisis pharetra at eget diam. Suspendisse ut vulputate magna. Nunc viverra posuere orci sit amet pulvinar. Quisque dui justo, egestas ac accumsan suscipit, tristique eu risus. In aliquet libero ante, ac pulvinar lectus pretium in. Ut enim tellus, vulputate et lorem et, hendrerit rutrum diam. Cras pharetra dapibus elit vitae ullamcorper. Nulla facilisis euismod diam, at sodales sem laoreet hendrerit. Curabitur porttitor ex nulla. Proin ligula leo, egestas ac nibh a, pellentesque mollis augue. Donec nec augue vehicula eros pulvinar vehicula eget rutrum neque. Duis sit amet interdum lorem. Vivamus porttitor nec quam a accumsan. Nam pretium vitae leo vitae rhoncus.</p>

                    <h2 id="drawer-jump-links-jump-links-fifth">Fifth section</h2>
                    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas sed dui ullamcorper, dignissim purus eu, mattis leo. Curabitur eleifend turpis ipsum, aliquam pretium risus efficitur vel. Etiam eget enim vitae quam facilisis pharetra at eget diam. Suspendisse ut vulputate magna. Nunc viverra posuere orci sit amet pulvinar. Quisque dui justo, egestas ac accumsan suscipit, tristique eu risus. In aliquet libero ante, ac pulvinar lectus pretium in. Ut enim tellus, vulputate et lorem et, hendrerit rutrum diam. Cras pharetra dapibus elit vitae ullamcorper. Nulla facilisis euismod diam, at sodales sem laoreet hendrerit. Curabitur porttitor ex nulla. Proin ligula leo, egestas ac nibh a, pellentesque mollis augue. Donec nec augue vehicula eros pulvinar vehicula eget rutrum neque. Duis sit amet interdum lorem. Vivamus porttitor nec quam a accumsan. Nam pretium vitae leo vitae rhoncus.</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="pf-c-drawer__panel" hidden="">
          <div class="pf-c-drawer__body">
            <div class="pf-c-drawer__head">
              <span>drawer-panel</span>
              <div class="pf-c-drawer__actions">
                <div class="pf-c-drawer__close">
                  <button class="pf-c-button pf-m-plain" type="button" aria-label="Close drawer panel">
                    <i class="fas fa-times" aria-hidden="true"></i>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>`)
}

func (h *Home) pageMain(uis ...app.UI) app.UI {
	el := app.Main().
		ID("main").
		Class("pf-c-page__main").
		TabIndex(-1)

	return el.Body(
		pageSession(
			func() app.UI {
				// list bordered

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
											Body(app.Text("Connect")).
											OnClick(func(ctx app.Context, e app.Event) {
												h.services = h.listServices(h.target)
												if len(h.services) > 0 {
													h.curSrv = h.services[0]
												}
												fmt.Println("services", h.services)
											})
									},
									func() app.UI {
										return app.Button().
											Class("pf-c-dropdown__toggle-button").
											Type("button").
											Aria("expanded", true).
											ID("dropdown-split-button-action-toggle-button").
											Aria("label", "Dropdown toggle").Body(
											app.I().
												Class("fas fa-caret-down").
												Aria("hidden", true),
										).OnClick(func(ctx app.Context, e app.Event) {
											h.saveExpanded = !h.saveExpanded
										})
									},
								))
							},

							// dropdown__menu
							func() app.UI {
								return app.Ul().
									Class("pf-c-dropdown__menu").
									Hidden(!h.saveExpanded).Body(jsutil.UIWrap(
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
											if len(h.methods) > 0 {
												h.curMth = h.methods[0]
											}
											fmt.Println(h.methods)
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
											schema, template := h.functionDescribe(h.target, h.curMth)
											h.input = template
											h.output = schema
											if h.editor == nil {
												h.editor = ace.New("editor")
												h.editor.SetReadOnly(true)
											}
											h.editor.SetValue(h.input)
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
							),
							app.Div().
								Class("pf-c-card__title").
								Text("Input"),
							app.Div().
								Class("pf-c-card__body").
								Body(

									app.Button().
										Class("pf-c-button pf-m-primary").
										Type("submit").
										Text("Submit").OnClick(func(ctx app.Context, e app.Event) {
										fmt.Println("close:", h.curMth)

										fmt.Println(h.editor.GetValue())
										var dd = h.invokeFunc(h.target, h.curMth, h.editor.GetValue())
										h.editor.SetValue(dd.Data.Result)
									}),
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
										return app.Div().
											Class("pf-c-code-editor").Body(jsutil.UIWrap(
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
								)),
						)
				}

				return app.Div().
					Class("pf-c-sidebar").
					Body(
						app.Div().
							Class("pf-c-sidebar__main").
							Body(
								app.Div().
									Class("pf-c-sidebar__panel pf-m-sticky").
									Body(
										pageSession(
											func() app.UI {
												return app.Div().
													Class("pf-c-search-input").
													Body(
														app.Div().
															Class("pf-c-search-input__bar").
															Body(
																app.Span().
																	Class("pf-c-search-input__text").
																	Body(
																		app.Span().
																			Class("pf-c-search-input__icon").
																			Body(
																				app.I().
																					Class("fas fa-search fa-fw").
																					Aria("hidden", true),
																			),
																		app.Input().Value(h.reqFilter).
																			Class("pf-c-search-input__text-input").
																			Type("text").
																			Placeholder("Find by name").
																			Aria("label", "Find by name").OnInput(func(ctx app.Context, e app.Event) {
																			h.reqFilter = ctx.JSSrc().Get("value").String()
																			fmt.Println(h.reqFilter)
																		}),
																	),
															),
													)
											},
											func() app.UI {
												return h.requestListUI()
											}),
									),
								app.Div().
									Class("pf-c-sidebar__content pf-m-no-background").
									Body(
										pageSession(func() app.UI {
											return cardForm()
										}),
									),
							),
					)

				return grid(
					gridItem(3)(h.requestListUI()),
					gridItem(4, 3)(cardForm()),
					gridItem(5, 7)(cardExample()),
					gridItem(6, 3)(cardEdit()),
				)
			}),
	)
}
func (h *Home) requestListUI() app.UI {
	el := app.Ul().Class("pf-c-list pf-m-plain pf-m-bordered")
	return el.Body(
		app.Range(h.allRequests).Slice(func(i int) app.UI {
			if !strings.Contains(strings.ToLower(h.allRequests[i].SelectedFunction), strings.ToLower(h.reqFilter)) {
				return nil
			}

			return app.Li().Class("pf-c-list__item").Body(
				app.A().Class("pf-c-list__item-icon").Body(
					app.I().Class("fas fa-times").
						Aria("hidden", true).OnClick(func(ctx app.Context, e app.Event) {
						fmt.Println("close:", h.allRequests[i].ID)
					}),
				),

				app.A().Class("pf-c-list__item-text").Body(
					app.Text(h.allRequests[i].Name)).OnClick(func(ctx app.Context, e app.Event) {
					fmt.Println("srv:", h.allRequests[i].ID)
					h.curSrv = h.allRequests[i].SelectedService
					h.curMth = h.allRequests[i].SelectedFunction
					h.target = h.allRequests[i].ServerTarget
					schema, template := h.functionDescribe(h.target, h.curMth)
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

func (h *Home) Render() app.UI {
	fmt.Println("Render", h.Mounted())
	return page(
		SkipContent(),
		pageHeader(),

		//h.pageSidebar(),

		h.pageMain(),
	)
}

func grid(items ...app.UI) app.UI {
	var el = app.Div().Class("pf-l-grid", "pf-m-gutter")
	return el.Body(items...)
}

func column(items ...app.UI) app.UI {
	var el = app.Div().Class("pf-l-grid", "pf-m-gutter", fmt.Sprintf("pf-m-all-%d-col", 12/len(items)))
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

func emptyItem() app.UI {
	return app.Div().Class("pf-l-grid__item")
}
