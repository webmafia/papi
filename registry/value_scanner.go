package registry

import (
	"reflect"

	"github.com/webmafia/fastapi/openapi"
	"github.com/webmafia/fastapi/registry/value"
)

type ValueScannerCreator interface {
	CreateScanner(tags reflect.StructTag) (scan value.ValueScanner, err error)
	Describe(schema *openapi.Schema) (err error)
}
