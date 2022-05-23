package jsutil

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var ErrNotFound = fmt.Errorf("query did not find element")

func QuerySelector(query string) (app.Value, error) {
	value := app.Window().Get("document").Call("querySelector", query)
	if !value.Truthy() {
		return app.Undefined(), ErrNotFound
	}
	return value, nil
}
