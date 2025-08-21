package registry

import (
	"reflect"
)

func (r *Registry) Binder(typ reflect.Type, fieldName string, tags reflect.StructTag) (_ Binder, err error) {

	// Use any existing binder
	if desc, ok := r.describe(typ); ok && desc.Binder != nil {
		return desc.Binder(fieldName, tags)
	}

	return
}
