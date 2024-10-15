package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webbmaffian/papi/registry/types"
	"github.com/webbmaffian/papi/route"
)

func (r *requestScanner) createParamScanner(typ reflect.Type, idx int) (scan types.RequestDecoder, err error) {
	sc, err := r.reg.CreateParamDecoder(typ, "")

	if err != nil {
		return
	}

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		params := route.RequestParams(c)
		return sc(p, params.Value(idx))
	}, nil
}
