package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
)

func (r *Registry) getCustomDecoder(typ reflect.Type, tags reflect.StructTag) (scan RequestDecoder, err error) {
	var dec Decoder

	// 1. If there is an explicit registered decoder, use it
	if desc, ok := r.desc[typ]; ok && desc.Decoder != nil {
		dec, err = desc.Decoder(tags)
	} else if typ.Implements(typeDescriber) {

		// 2. If the type can describe itself, let it
		if v, ok := reflect.New(typ).Interface().(TypeDescriber); ok {
			if desc := v.TypeDescription(r); desc.Schema != nil {
				dec, err = desc.Decoder(tags)
			}
		}
	}

	if err == nil && dec != nil {
		scan = func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
			return dec(p, fast.BytesToString(c.Request.Body()))
		}
	}

	return
}
