package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
)

func (r *Registry) createQueryDecoder(typ reflect.Type, key string, tags reflect.StructTag) (scan RequestDecoder, err error) {
	sc, err := r.Decoder(typ, tags)

	if err != nil {
		return
	}

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		val := c.QueryArgs().Peek(key)

		if len(val) > 0 {
			if err := sc(p, fast.BytesToString(val)); err != nil {
				return ErrFailedDecodeQuery.Detailed(err.Error(), key)
			}
		}

		return nil
	}, nil
}
