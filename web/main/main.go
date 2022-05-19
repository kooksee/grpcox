package main

import (
	_ "github.com/gusaul/grpcox/web/app"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	app.RunWhenOnBrowser()
}
