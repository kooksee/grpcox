package app

import (
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/maxence-charriere/go-app/v9/pkg/ui"
)

const (
	backgroundColor = "#000000"
)

func init() {
	ui.BaseHPadding = 42
	ui.BlockPadding = 18
}

func App() *app.Handler {
	return &app.Handler{
		Name:            "turn gRPCurl into web based UI, extremely easy to use",
		Title:           "grpc-ui",
		Description:     "An Home World! example",
		Author:          "Maxence Charriere",
		Image:           "https://go-app.dev/web/images/go-app.png",
		BackgroundColor: backgroundColor,
		//BackgroundColor: "#111",
		ThemeColor:   backgroundColor,
		LoadingLabel: "go-app documentation",
		Body: func() app.HTMLBody {
			return app.Body()
		},
		Scripts: []string{
			//"/js/jquery-3.3.1.min.js",
			//"/js/popper.min.js",
			//"/js/bootstrap.bundle.min.js",
			//"/js/bootstrap.min.js",
			//"/js/mdb.min.js",
			//"http://cdn.staticfile.org/prettify/r298/prettify.min.js",
			//"/js/ace.js",
			//"https://cdnjs.cloudflare.com/ajax/libs/material-components-web/13.0.0/material-components-web.js",
			//"https://cdnjs.cloudflare.com/ajax/libs/prism/1.25.0/prism.min.js",
			//"https://cdnjs.cloudflare.com/ajax/libs/prism/1.25.0/components/prism-go.min.js",

			//"/js/style.js",
			//"/js/db.js",
			//"/js/proto.js",
			//"/js/ctx.metadata.js",
			//"/js/request.list.js",
		},
		Styles: []string{
			"https://unpkg.com/@patternfly/patternfly@4.135.2/patternfly.css",
			"https://unpkg.com/@patternfly/patternfly@4.135.2/patternfly-addons.css",
			//"https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css",
			//"/css/bootstrap.min.css",
			//"/css/mdb.min.css",
			//"/css/style.css",
			//"/css/proto.css",
			//"https://fonts.googleapis.com/icon?family=Material+Icons",
			//"https://fonts.googleapis.com/css2?family=Roboto&display=swap",
			//"https://cdnjs.cloudflare.com/ajax/libs/material-components-web/13.0.0/material-components-web.css",
			//"https://cdnjs.cloudflare.com/ajax/libs/prism-themes/1.9.0/prism-material-light.min.css",
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
