package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/internal/json"
)

func (r *Registry) createJsonDecoder(typ reflect.Type) (scan Handler, err error) {
	dec := json.DecoderOf(typ)
	scan = func(c *fasthttp.RequestCtx, p unsafe.Pointer) error {
		iter := json.AcquireIterator(c.Request.BodyStream())
		defer json.ReleaseIterator(iter)

		dec.Decode(p, iter)

		if iter.Error != nil {
			return ErrFailedDecodeBody.Detailed(iter.Error.Error())
		}

		return nil
	}

	return
}
