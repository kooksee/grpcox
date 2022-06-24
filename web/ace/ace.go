package ace

import (
	"fmt"

	"github.com/gusaul/grpcox/web/jsutil"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// https://ace.c9.io/#nav=howto&api=ace

func New(id string) *Ace {
	CfgSet("basePath", "https://ace.c9.io/build/src")

	var ace = &Ace{val: jsutil.NewAtPath("ace.edit", id, map[string]interface{}{
		"theme":                    "ace/theme/github",
		"mode":                     "ace/mode/json",
		"autoScrollEditorIntoView": true,
		"maxLines":                 30,
		"minLines":                 2,
	})}
	ace.SetShowGutter(false)
	ace.SetScrollMargin(10, 10, 10, 10)

	return ace
}

type Ace struct {
	val app.Value
}

func SetModuleUrl(key string, val string) {
	jsutil.Eval(fmt.Sprintf("ace.config.setModuleUrl('%s', '%s')", key, val))
}

func CfgSet(key string, val string) {
	jsutil.Eval(fmt.Sprintf("ace.config.set('%s', '%s')", key, val))
}

// SetTheme editor.setTheme("ace/theme/monokai")
// ace/theme/tomorrow_night_blue
// ace/theme/github
func (t *Ace) SetTheme(val string) {
	t.val.Call("setTheme", val)
}

func (t *Ace) SetReadOnly(val bool) {
	t.val.Call("setReadOnly", val)
}

// SetMode editor.session.setMode("ace/mode/javascript");
// ace/mode/html
// ace/mode/json
func (t *Ace) SetMode(val string) {
	t.val.Get("session").Call("setMode", val)
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

// SetShowGutter "SetShowGutter"
func (t *Ace) SetShowGutter(args ...interface{}) {
	t.val.Get("renderer").Call("setShowGutter", args...)
}

func (t *Ace) GetValue() string {
	return t.val.Call("getValue").String()
}
