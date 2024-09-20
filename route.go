package fastapi

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/webmafia/fastapi/internal/jsonpool"
)

func (api *API) RegisterRoutes(types ...any) (err error) {
	for i := range types {
		val := reflect.ValueOf(types[i])
		numMethods := val.NumMethod()

		for i := 0; i < numMethods; i++ {
			cb, ok := val.Method(i).Interface().(func(api *API) error)

			if !ok {
				return errors.New("invalid handler")
			}

			if err = cb(api); err != nil {
				return
			}
		}
	}

	return
}

func AddRoute[I, O any](api *API, r Route[I, O]) (err error) {
	iTyp := reflect.TypeOf((*I)(nil)).Elem()
	oTyp := reflect.TypeOf((*O)(nil)).Elem()
	_ = oTyp
	route := api.router.Add(string(r.Method), r.Path)

	cb, err := createInputScanner(iTyp, route.params)

	if err != nil {
		return
	}

	route.handler = func(ctx *Ctx) (err error) {
		ctx.ctx.SetContentType("application/json; charset=utf-8")

		s := jsonpool.AcquireStream(ctx.ctx.Response.BodyWriter())
		defer jsonpool.ReleaseStream(s)

		var (
			in     I
			out    O
			outAny any = &out
		)

		if err = cb(unsafe.Pointer(&in), ctx.ctx, ctx.paramVals); err != nil {
			return
		}

		if enc, ok := outAny.(Lister); ok {
			s.WriteObjectStart()
			s.WriteObjectField("items")
			s.WriteArrayStart()

			enc.setStream(s)

			if err = r.Handler(ctx, &in, &out); err != nil {
				return
			}

			s.WriteArrayEnd()
			s.WriteMore()

			s.WriteObjectField("meta")
			enc.encodeMeta(s)

			s.WriteObjectEnd()
		} else {
			if err = r.Handler(ctx, &in, &out); err != nil {
				return
			}

			if enc, ok := outAny.(JsonEncoder); ok {
				if err = enc.EncodeJson(s); err != nil {
					return
				}
			} else {
				s.WriteVal(out)
			}
		}

		return s.Flush()
	}

	return
}
