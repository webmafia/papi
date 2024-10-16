package registry

import (
	"reflect"

	"github.com/webbmaffian/papi/openapi"
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

	// Request & response
	Schema func(tags reflect.StructTag) (openapi.Schema, error)

	// Request
	Handler func(tags reflect.StructTag, handler Handler) (Handler, error)

	// Response
	Decoder func(tags reflect.StructTag) (Decoder, error)
}

func (t TypeDescription) IsZero() bool {
	return t.Schema == nil && t.Handler == nil && t.Decoder == nil
}
