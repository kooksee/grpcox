package jsutil

import (
	"fmt"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
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

func IfElse(ok bool, a, b interface{}) interface{} {
	if ok {
		return a
	}

	return b
}

func UIWrap(uis ...func() app.UI) app.UI {
	var uiList []app.UI
	for i := range uis {
		uiList = append(uiList, uis[i]())
	}
	return app.If(true, uiList...)
}

type UI []func() app.UI

func (u UI) Render() app.UI {
	var uiList []app.UI
	for i := range u {
		uiList = append(uiList, u[i]())
	}
	return app.If(true, uiList...)
}

type Classes map[string]bool

func (t Classes) Render() string {
	var data []string
	for name := range t {
		if !t[name] {
			continue
		}
		data = append(data, name)
	}
	return strings.Join(data, " ")
}
