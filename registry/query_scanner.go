package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webbmaffian/papi/registry/types"
	"github.com/webmafia/fast"
)

func (r *requestScanner) createQueryScanner(typ reflect.Type, key string) (scan types.RequestDecoder, err error) {
	sc, err := r.reg.CreateParamDecoder(typ, "")

	if err != nil {
		return
	}

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		val := c.QueryArgs().Peek(key)

		if len(val) > 0 {
			return sc(p, fast.BytesToString(val))
		}

		return nil
	}, nil
}
