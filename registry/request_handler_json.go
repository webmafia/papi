package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/internal/json"
)

func (r *Registry) createJsonBinder(typ reflect.Type) (bind Binder, err error) {
	dec := json.DecoderOf(typ)
	bind = func(c *fasthttp.RequestCtx, p unsafe.Pointer) error {
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
