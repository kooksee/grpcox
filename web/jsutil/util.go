package jsutil

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"strings"
)

func Log(args ...interface{}) {
	app.Window().Get("console").Call("log", args...)
}

type ValueHelper struct {
	Root app.Value
}

func (v ValueHelper) Get(path ...string) (out app.Value, ok bool) {
	for _, part := range path {
		next := v.Root.Get(part)
		if next.IsUndefined() {
			return nil, false
		}
		v.Root = next
	}
	if v.Root.IsUndefined() {
		return nil, false
	}
	return v.Root, true
}

func (v ValueHelper) List(path ...string) (out []app.Value, ok bool) {
	if list, ok := v.Get(path...); ok {
		length := list.Length()
		for i := 0; i < length; i++ {
			out = append(out, list.Index(i))
		}
		return out, true
	}
	return nil, false
}

func Eval(code string) app.Value {
	return app.Window().Call("eval", code)
}

func ValueAtPath(path string) app.Value {
	return ValueAt(app.Window(), path, true)
}

func ValueAt(root app.Value, path string, warn bool) app.Value {
	if root == nil || root.IsUndefined() {
		return app.Undefined()
	}
	var current = root
	for _, part := range strings.Split(path, ".") {
		current = current.Get(part)
		if current.IsUndefined() {
			if warn {
				fmt.Printf("jsutil.getValue(%q) : undefined at %q\n", path, part)
			}
			break
		}
	}
	return current
}

func NewAtPath(path string, args ...interface{}) app.Value {
	v := ValueAtPath(path)
	if v.Truthy() {
		return v.New(args...)
	}
	return app.Undefined()
}

func CallAtPath(path string, m string, args ...interface{}) app.Value {
	v := ValueAtPath(path)
	if v.Truthy() {
		return v.Call(m, args...)
	}
	return app.Undefined()
}

func Hidden(ok bool) interface{} {
	if ok {
		return true
	}
	return ""
}

func If(ok bool, val interface{}) interface{} {
	if ok {
		return val
	}

	return nil
}

type Compo struct {
	uiList []app.UI
}

func (c *Compo) UI(args ...app.UI) {
	c.uiList = append(c.uiList, args...)
}

func (c *Compo) Fn(args ...func() app.UI) {
	var uiList []app.UI
	for i := range args {
		uiList = append(uiList, args[i]())
	}
	c.uiList = append(c.uiList, uiList...)
}

func (c *Compo) Body(c1 func(c *Compo)) app.UI { return Body(c1) }

func Body(c func(c *Compo)) app.UI {
	var cc = &Compo{}
	c(cc)
	return app.If(true, cc.uiList...)
}

func UI(name string, args ...func() app.UI) app.UI {
	var uiList []app.UI
	for i := range args {
		uiList = append(uiList, args[i]())
	}
	return app.If(true, uiList...)
}

func Wrap(name string, args ...func() app.UI) func() app.UI {
	return func() app.UI {
		var uiList []app.UI
		for i := range args {
			uiList = append(uiList, args[i]())
		}
		return app.If(true, uiList...)
	}
}

type Component struct {
	app.Compo
}

func (c Component) init() {

}
