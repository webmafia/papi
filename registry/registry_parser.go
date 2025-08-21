package registry

import (
	"reflect"
)

// Get a parser, or create one.
func (r *Registry) Parser(typ reflect.Type, tags reflect.StructTag) (dec Parser, err error) {

	// Use any existing parser
	if desc, ok := r.describe(typ); ok && desc.Parser != nil {
		return desc.Parser(tags)
	}

	// Generate a parser
	return r.scan.CreateScanner(typ)
}
