package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/internal/route"
)

func (r *Registry) createParamBinder(typ reflect.Type, key string, idx int, tags reflect.StructTag) (_ Binder, err error) {
	sc, err := r.Parser(typ, tags)

	if err != nil {
		return
	}

	return func(c *fasthttp.RequestCtx, p unsafe.Pointer) error {
		params := route.RequestParams(c)

		if err := sc(p, params.Value(idx)); err != nil {
			return ErrFailedDecodeParam.Detailed(err.Error(), key)
		}

		return nil
	}, nil
}
