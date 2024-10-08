package types

import (
	"reflect"

	"github.com/webmafia/fastapi/openapi"
	"github.com/webmafia/fastapi/registry/scanner"
)

type Type interface {
	Type() reflect.Type
	CreateScanner(tags reflect.StructTag) (scan scanner.Scanner, err error)
	Describe(tags reflect.StructTag) (schema openapi.Schema)
}
