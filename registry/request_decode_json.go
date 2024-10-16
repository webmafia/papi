package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
)

func (r *Registry) createJsonDecoder(typ reflect.Type) (scan RequestDecoder, err error) {
	dec := r.json.DecoderOf(typ)
	scan = func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		iter := r.json.AcquireIterator(c.Request.BodyStream())
		defer r.json.ReleaseIterator(iter)

		dec.Decode(p, iter)
		return iter.Error
	}

	return
}
