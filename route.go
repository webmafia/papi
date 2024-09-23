package fastapi

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
)

func (api *API) RegisterRoutes(types ...any) (err error) {
	for i := range types {
		val := reflect.ValueOf(types[i])
		numMethods := val.NumMethod()

		for i := 0; i < numMethods; i++ {
			cb, ok := val.Method(i).Interface().(func(api *API) error)

			if !ok {
				continue
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

	cb, err := createInputScanner(api, iTyp, route.params)

	if err != nil {
		return
	}

	route.handler = func(c *fasthttp.RequestCtx) (err error) {
		c.SetContentType("application/json; charset=utf-8")

		s := api.opt.JsonPool.AcquireStream(c.Response.BodyWriter())
		defer api.opt.JsonPool.ReleaseStream(s)

		var (
			in     I
			out    O
			outAny any = &out
		)

		if err = cb(unsafe.Pointer(&in), c); err != nil {
			return
		}

		if enc, ok := outAny.(Lister); ok {
			s.WriteObjectStart()
			s.WriteObjectField("items")
			s.WriteArrayStart()

			enc.setStream(s)

			if err = r.Handler(c, &in, &out); err != nil {
				return
			}

			s.WriteArrayEnd()
			s.WriteMore()

			s.WriteObjectField("meta")
			enc.encodeMeta(s)

			s.WriteObjectEnd()
		} else {
			if err = r.Handler(c, &in, &out); err != nil {
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
