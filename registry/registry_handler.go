package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
)

type Handler func(c *fasthttp.RequestCtx, in, out unsafe.Pointer) error

func (r *Registry) Handler(typ reflect.Type, tags reflect.StructTag, paramKeys []string, handler Handler) (scan Handler, err error) {

	// 1. If there is an explicit registered handler describer, use it
	if desc, ok := r.desc[typ]; ok {
		return desc.Handler(tags, handler)
	}

	// 2. If the type can describe itself, let it
	if typ.Implements(typeDescriber) {
		if v, ok := reflect.New(typ).Interface().(TypeDescriber); ok {
			if desc := v.TypeDescription(r); desc.Schema != nil {
				return desc.Handler(tags, handler)
			}
		}
	}

	// 3. In all other cases, use the default handler
	return r.defaultHandler(typ, handler)
}

func (r *Registry) defaultHandler(typ reflect.Type, handler Handler) (Handler, error) {
	enc := r.json.EncoderOf(typ)

	return func(c *fasthttp.RequestCtx, in, out unsafe.Pointer) error {
		if err := handler(c, in, out); err != nil {
			return err
		}

		c.SetContentType("application/json")

		s := r.json.AcquireStream(c.Response.BodyWriter())
		defer r.json.ReleaseStream(s)

		enc.Encode(out, s)
		return s.Flush()
	}, nil
}
