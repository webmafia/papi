package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/internal/json"
)

type Handler func(c *fasthttp.RequestCtx, ptr unsafe.Pointer) error

func (r *Registry) Handler(typ reflect.Type, handler Handler) (scan Handler, err error) {

	// 1. If there is an explicit registered handler describer, use it
	if desc, ok := r.desc[typ]; ok {
		return desc.Handler(handler)
	}

	// 2. If the type can describe itself, let it
	if typ.Implements(typeDescriber) {
		if v, ok := reflect.New(typ).Interface().(TypeDescriber); ok {
			if desc := v.TypeDescription(r); desc.Handler != nil {
				return desc.Handler(handler)
			}
		}
	}

	// 3. In all other cases, use the default handler
	return r.defaultHandler(typ, handler)
}

func (r *Registry) defaultHandler(typ reflect.Type, handler Handler) (Handler, error) {
	enc := json.EncoderOf(typ)

	return func(c *fasthttp.RequestCtx, ptr unsafe.Pointer) error {
		if err := handler(c, ptr); err != nil {
			return err
		}

		c.SetContentType("application/json")

		s := json.AcquireStream(c.Response.BodyWriter())
		defer json.ReleaseStream(s)

		enc.Encode(ptr, s)
		return s.Flush()
	}, nil
}
