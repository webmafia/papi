package datatype

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
)

type fieldScanner struct {
	scan   func(ctx *fasthttp.RequestCtx, p unsafe.Pointer, params []string) error
	offset uintptr
}

func CreateRequestScanner[T any](d *DataTypes) (fn func(ctx *fasthttp.RequestCtx, p *T, params []string) error, err error) {
	typ := reflect.TypeOf((*T)(nil)).Elem()

	if kind := typ.Kind(); kind != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", kind)
	}

	numFields := typ.NumField()
	scans := make([]fieldScanner, 0, numFields)

	for i := 0; i < numFields; i++ {
		fld := typ.Field(i)

		if fld.Name == "Body" {
			fldTyp := fld.Type

			if fldTyp.Kind() != reflect.Pointer {
				fldTyp = reflect.PointerTo(fldTyp)
			}

			// Todo: Check for body scanner, with JSON fallback for structs
		} else {

		}
	}
}
