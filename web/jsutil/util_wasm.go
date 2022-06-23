package jsutil

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pubgo/xerror"
	"honnef.co/go/js/dom/v2"
	fetch "marwan.io/wasm-fetch"
)

func Event(value app.Value) dom.Event {
	return dom.WrapEvent(app.JSValue(value))
}

func LoadJs(jsUrl string) app.Value {
	rsp, err := fetch.Fetch(jsUrl, &fetch.Opts{Method: fetch.MethodGet})
	xerror.Panic(err)
	return Eval(string(rsp.Body))
}
