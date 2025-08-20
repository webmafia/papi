package registry

import (
	"reflect"

	"github.com/webmafia/papi/openapi"
)

var typeDescriber = reflect.TypeFor[TypeDescriber]()

type TypeDescriber interface {
	TypeDescription(reg *Registry) TypeDescription
}

type TypeRegistrar interface {
	Type() reflect.Type
	TypeDescriber
}

type TypeDescription struct {

	// Documentation of request and response
	Schema func(tags reflect.StructTag) (openapi.Schema, error)

	// Handler of request and response
	Handler func(handler Handler) (Handler, error)

	// Decoding of request value (e.g. query param)
	Decoder func(tags reflect.StructTag) (Decoder, error)
}

func (t TypeDescription) IsZero() bool {
	return t.Schema == nil && t.Handler == nil && t.Decoder == nil
}
