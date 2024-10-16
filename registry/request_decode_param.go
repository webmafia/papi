package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webbmaffian/papi/route"
)

func (r *Registry) createParamDecoder(typ reflect.Type, idx int, tags reflect.StructTag) (scan RequestDecoder, err error) {
	sc, err := r.Decoder(typ, tags)

	if err != nil {
		return
	}

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		params := route.RequestParams(c)
		return sc(p, params.Value(idx))
	}, nil
}
