package ui

import (
	"fmt"
	"reflect"
)

func (ctx *Context) Reflect(name string, data interface{}) {
	Reflect(ctx, name, data)
}

func Reflect(ctx *Context, name string, data interface{}) {
	rowHeight := ctx.Measure("Õj").Y
	indent := rowHeight

	var walk func(string, reflect.Value)
	walk = func(name string, v reflect.Value) {
		// walk pointers
		for v.Kind() == reflect.Ptr {
			if v.IsNil() {
				ctx.Top(rowHeight).Text(fmt.Sprintf("%v = nil", name))
			}
			v = v.Elem()
		}

		if v.Kind() == reflect.Struct {
			ctx.Top(rowHeight).Text(name)

			// recurse
			t := v.Type()
			for i, n := 0, v.NumField(); i < n; i++ {
				f, ft := v.Field(i), t.Field(i)

				ctx.Area.Min.X += indent
				walk(ft.Name, f)
				ctx.Area.Min.X -= indent
			}
		} else {
			ctx.Top(rowHeight).Text(fmt.Sprintf("%v = %v", name, v.Interface()))
		}
	}
	walk(name, reflect.ValueOf(data))
}
