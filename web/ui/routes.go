package app

import "github.com/maxence-charriere/go-app/v9/pkg/app"

func init() {
	app.Route("/", &Home{})
}
