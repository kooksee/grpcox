package jsutil

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"honnef.co/go/js/dom/v2"
)

func Event(value app.Value) dom.Event {
	return dom.WrapEvent(app.JSValue(value))
}
