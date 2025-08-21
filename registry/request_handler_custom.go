package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
)

// Returns a nil handler if there is no custom handler.
func (r *Registry) getCustomBinder(typ reflect.Type, tags reflect.StructTag) (bind Binder, err error) {
	var dec Parser

	// Use any existing binder
	if desc, ok := r.describe(typ); ok && desc.Binder != nil {
		return desc.Binder(tags)
	}

	if err == nil && dec != nil {
		bind = func(c *fasthttp.RequestCtx, p unsafe.Pointer) error {
			return dec(p, fast.BytesToString(c.Request.Body()))
		}
	}

	return
}
