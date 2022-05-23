package jsutil

import "github.com/maxence-charriere/go-app/v9/pkg/app"

func Alert(message string) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		e.PreventDefault()
		app.Window().Call("alert", message)
	}
}
