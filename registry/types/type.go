package types

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webbmaffian/papi/openapi"
	"github.com/webbmaffian/papi/registry/scanner"
)

type Type interface {
	Type() reflect.Type
}

type ParamType interface {
	Type
	ParamSchema(tags reflect.StructTag) (schema openapi.Schema)
	CreateParamScanner(tags reflect.StructTag) (scan scanner.Scanner, err error)
}

type RequestType interface {
	Type
	CreateRequestScanner(tags reflect.StructTag, paramKeys []string) (scan RequestScanner, err error)
	DescribeOperation(op *openapi.Operation) (err error)
}

type RequestScanner func(p unsafe.Pointer, c *fasthttp.RequestCtx) error
