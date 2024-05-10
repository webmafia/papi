package fastapi

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/webmafia/fastapi/internal/jsonpool"
	"github.com/webmafia/fastapi/scan"
)

func (api *API[U]) RegisterRoutes(types ...any) (err error) {
	for i := range types {
		val := reflect.ValueOf(types[i])
		numMethods := val.NumMethod()

		for i := 0; i < numMethods; i++ {
			cb, ok := val.Method(i).Interface().(func(api *API[U]) error)

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

func AddRoute[U, I, O any](api *API[U], r Route[U, I, O]) (err error) {
	var in I
	inTyp := reflect.TypeOf(in)
	route := api.router.Add(string(r.Method), r.Path)
	sc, err := scan.CreateStructScanner(inTyp, route.params)

	if err != nil {
		return
	}

	if err = addRouteDocs(api, r); err != nil {
		return
	}

	route.cb = func(ctx *Ctx[U]) (err error) {
		ctx.ctx.SetContentType("application/json; charset=utf-8")

		s := jsonpool.AcquireStream(ctx.ctx.Response.BodyWriter())
		defer jsonpool.ReleaseStream(s)

		var (
			in     I
			out    O
			outAny any = &out
		)

		if err = sc(unsafe.Pointer(&in), ctx.ctx, ctx.paramVals); err != nil {
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
