package ace

import (
	"github.com/gusaul/grpcox/web/jsutil"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func New(id string) *Ace {
	return &Ace{val: jsutil.NewAtPath("ace.edit", id)}
}

type Ace struct {
	val app.Value
}

func (t *Ace) SetValue(val string) {
	t.val.Call("setValue", val)
}

func (t *Ace) SetOptions(data map[string]interface{}) {
	t.val.Call("setOptions", data)
}

func (t *Ace) SetScrollMargin(args ...interface{}) {
	t.val.Get("renderer").Call("setScrollMargin", args...)
}

func (t *Ace) getValue() string {
	return t.val.Call("getValue").String()
}
