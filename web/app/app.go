package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pubgo/xerror"
)

const (
	backgroundColor = "#000000"
)

func App() *app.Handler {
	//ui.BaseHPadding = 42
	//ui.BlockPadding = 18

	return &app.Handler{
		Name:            "turn gRPCurl into web based UI, extremely easy to use",
		Title:           "grpc-ui",
		Description:     "An Home World! example",
		Author:          "Maxence Charriere",
		Image:           "https://go-app.dev/web/images/go-app.png",
		BackgroundColor: backgroundColor,
		ThemeColor:      backgroundColor,
		LoadingLabel:    "go-app documentation",
		Body: func() app.HTMLBody {
			return app.Body()
		},
		Scripts: []string{
			"/js/jquery-3.3.1.min.js",
			"/js/popper.min.js",
			"/js/bootstrap.min.js",
			"/js/mdb.min.js",
			"https://cdnjs.cloudflare.com/ajax/libs/prettify/r298/run_prettify.js",
			"/js/ace.js",
			"/js/db.js",
			"/js/style.js",
			"/js/proto.js",
			"/js/ctx.metadata.js",
			"/js/request.list.js",
		},
		Styles: []string{
			"https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css",
			"/css/bootstrap.min.css",
			"/css/mdb.min.css",
			"/css/style.css",
			"/css/proto.css",
		},
		RawHeaders: []string{
			`<link rel="icon" href="img/favicon.png" type="image/x-icon" />`,
			`<meta charset="utf-8">`,
			`<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">`,
			`<meta http-equiv="x-ua-compatible" content="ie=edge">`,
			`<title>gRPCox - gRPC Testing Environment</title>`,
		},
		CacheableResources: []string{
			"/img/favicon.png",
		},
		AutoUpdateInterval: time.Second * 5,
	}
}

type Requests struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	RawRequest       string `json:"raw_request"`
	ResponseHTML     string `json:"response_html"`
	SchemaProtoHTML  string `json:"schema_proto_html"`
	SelectedFunction string `json:"selected_function"`
	SelectedService  string `json:"selected_service"`
	ServerTarget     string `json:"server_target"`
}

type Home struct {
	app.Compo
	data map[string]*Requests
}

func (h *Home) OnInit() {
	h.data = make(map[string]*Requests)
}

func (h *Home) OnMount(ctx app.Context) {
	ctx.Async(func() {
		rsp, err := http.Get("/api/requests")
		xerror.Panic(err)

		var req map[string][]*Requests
		dtBytes, err := ioutil.ReadAll(rsp.Body)
		xerror.Panic(err)
		xerror.Panic(json.Unmarshal(dtBytes, &req))

		for _, r := range req["data"] {
			h.data[r.ID] = r
		}

		for k := range h.data {
			fmt.Println(k)
		}
		h.Update()
	})
}
func (h *Home) OnNav(ctx app.Context) {
	fmt.Println("component navigated")
}

func (h *Home) Render() app.UI {
	return app.Div().Body(
		app.Raw(`
<button type="button" class="btn btn-primary" data-toggle="modal" data-target="#basicExampleModal">
  Launch demo modal
</button>
`),

		app.Raw(`
<div class="modal fade" id="basicExampleModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalLabel"
  aria-hidden="true">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="exampleModalLabel">Modal title</h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        ...
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
        <button type="button" class="btn btn-primary">Save changes</button>
      </div>
    </div>
  </div>
</div>
`),
	)
	//return app.Div().Body(
	//	app.Range(h.data).Map(func(s string) app.UI {
	//		return app.Div().Text(h.data[s].Name)
	//	}),
	//)
}
