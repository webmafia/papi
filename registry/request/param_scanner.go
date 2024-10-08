package request

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webbmaffian/papi/registry"
	"github.com/webbmaffian/papi/route"
)

const paramsKey = "params"

func (r *requestScanner) createParamScanner(typ reflect.Type, idx int) (scan registry.RequestScanner, err error) {
	sc, err := r.reg.CreateValueScanner(typ, "")

	if err != nil {
		return
	}

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		params := RequestParams(c)
		return sc(p, params.Value(idx))
	}, nil
}

func RequestParams(c *fasthttp.RequestCtx) *route.Params {
	if params, ok := c.UserValue(paramsKey).(*route.Params); ok {
		return params
	}

	return route.NilParams
}

func SetRequestParams(c *fasthttp.RequestCtx, params *route.Params) {
	c.SetUserValue(paramsKey, params)
}
