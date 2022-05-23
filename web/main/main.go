package main

import (
	_ "github.com/gusaul/grpcox/web/ui"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pubgo/xerror"
)

func main() {
	for {
		_ = xerror.Try(func() {
			defer xerror.RespDebug()
			app.RunWhenOnBrowser()
		})
	}
}
