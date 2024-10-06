package request

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
	"github.com/webmafia/fastapi/scanner"
)

func (r *requestScanner) createQueryScanner(typ reflect.Type, key string) (scan scanner.RequestScanner, err error) {
	sc, err := r.reg.CreateValueScanner(typ, "")

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
