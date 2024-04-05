package fastapi

import (
	"errors"
	"reflect"
	"unsafe"
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

	cb := func(ctx *Ctx[U]) error {
		var (
			in  I
			out O
		)

		// TODO: Create input scanner and handle output.

		return r.Handler(ctx, &in, &out)
	}

	api.router.Add(string(r.Method), r.Path, unsafe.Pointer(&cb))
	return
}

func parseIn[U any, I any](ctx *Ctx[U], in *I) (err error) {
	return
}
