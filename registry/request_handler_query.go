package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
)

func (r *Registry) createQueryBinder(typ reflect.Type, key string, tags reflect.StructTag) (_ Binder, err error) {
	sc, err := r.Parser(typ, tags)

	if err != nil {
		return
	}

	return func(c *fasthttp.RequestCtx, p unsafe.Pointer) error {
		val := c.QueryArgs().Peek(key)

		if len(val) > 0 {
			if err := sc(p, fast.BytesToString(val)); err != nil {
				return ErrFailedDecodeQuery.Detailed(err.Error(), key)
			}
		}

		return nil
	}, nil
}
