package fastapi

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/webmafia/fastapi/scan"
)

func (api *API[U]) RegisterRoutes(types ...any) (err error) {
	for i := range types {
		val := reflect.ValueOf(types[i])
		numMethods := val.NumMethod()

		for i := 0; i < numMethods; i++ {
			cb, ok := val.Method(i).Interface().(func(api *API[U]))

			if !ok {
				return errors.New("invalid handler")
			}

			cb(api)
		}
	}

	return
}

func AddRoute[U, I, O any](api *API[U], r Route[U, I, O]) (err error) {
	var in I
	sc, err := scan.CreateStructScanner(reflect.TypeOf(in))

	if err != nil {
		return
	}

	// TODO: Add to struct scanner

	scanStruct := sc.Compile()

	cb := func(ctx *Ctx[U]) (err error) {
		var (
			in  I
			out O
		)

		if err = scanStruct(unsafe.Pointer(&in)); err != nil {
			return
		}

		// TODO: Create input scanner and handle output.

		return r.Handler(ctx, &in, &out)
	}

	api.router.Add(string(r.Method), r.Path, cb)
	return
}

func parseIn[U any, I any](ctx *Ctx[U], in *I) (err error) {
	return
}
