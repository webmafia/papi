package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/internal/route"
)

func (r *Registry) createParamDecoder(typ reflect.Type, key string, idx int, tags reflect.StructTag) (scan RequestDecoder, err error) {
	sc, err := r.Decoder(typ, tags)

	if err != nil {
		return
	}

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		params := route.RequestParams(c)

		if err := sc(p, params.Value(idx)); err != nil {
			return ErrFailedDecodeParam.Detailed(err.Error(), key)
		}

		return nil
	}, nil
}
