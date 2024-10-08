package types

import (
	"reflect"

	"github.com/webmafia/fastapi/openapi"
	"github.com/webmafia/fastapi/registry/value"
)

type Type interface {
	Type() reflect.Type
	CreateScanner(tags reflect.StructTag) (scan value.ValueScanner, err error)
	Describe(tags reflect.StructTag) (schema *openapi.Schema)
}
