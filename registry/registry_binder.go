package registry

import (
	"reflect"
)

func (r *Registry) Binder(typ reflect.Type, tags reflect.StructTag) (scan Binder, err error) {

	// Use any existing binder
	if desc, ok := r.describe(typ); ok && desc.Binder != nil {
		return desc.Binder(tags)
	}

	return
}
