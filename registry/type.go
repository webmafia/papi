package registry

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
	CreateParamDecoder(tags reflect.StructTag) (scan ParamDecoder, err error)
}

type RequestType interface {
	Type
	CreateRequestDecoder(tags reflect.StructTag, paramKeys []string) (RequestDecoder, error)
	CreateResponseEncoder(reg *Registry, tags reflect.StructTag, paramKeys []string, handler ResponseEncoder) (ResponseEncoder, error)
	DescribeOperation(op *openapi.Operation) (err error)
}

type RequestDecoder func(p unsafe.Pointer, c *fasthttp.RequestCtx) error
type ResponseEncoder func(c *fasthttp.RequestCtx, in, out unsafe.Pointer) error
type ParamDecoder = scanner.Scanner
