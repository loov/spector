package ui

import (
	"fmt"
	"reflect"
)

func (ctx *Context) Reflect(data interface{}) { Reflect(ctx, data) }

func Reflect(ctx *Context, data interface{}) {
	rowHeight := ctx.Measure("Õj").Y

	ctx.Top(rowHeight / 2)

	names := ctx.Left(ctx.Area.Dx() / 3)
	values := ctx.Fill()

	v := reflect.ValueOf(data)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := reflect.TypeOf(data)

	if v.Kind() == reflect.Struct {
		for i, n := 0, v.NumField(); i < n; i++ {
			f := v.Field(i)
			ft := t.Field(i)
			names.Top(rowHeight).Text(ft.Name)
			values.Top(rowHeight).Text(fmt.Sprintf("%v", f.Interface()))
		}
	}
}
