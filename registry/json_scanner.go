package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webbmaffian/papi/registry/types"
)

func (r *requestScanner) createJsonScanner(typ reflect.Type) (scan types.RequestScanner, err error) {
	dec := r.json.DecoderOf(typ)
	scan = func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		iter := r.json.AcquireIterator(c.Request.BodyStream())
		defer r.json.ReleaseIterator(iter)

		dec.Decode(p, iter)
		return iter.Error
	}

	return
}
