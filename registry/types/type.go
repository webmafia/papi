package types

import (
	"reflect"

	"github.com/webbmaffian/papi/openapi"
	"github.com/webbmaffian/papi/registry/scanner"
)

type Type interface {
	Type() reflect.Type
	CreateScanner(tags reflect.StructTag) (scan scanner.Scanner, err error)
	Describe(tags reflect.StructTag) (schema openapi.Schema)
}
