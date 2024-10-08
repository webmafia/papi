package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webbmaffian/papi/openapi"
)

type (
	RequestScanner        func(p unsafe.Pointer, c *fasthttp.RequestCtx) error
	RequestScannerCreator interface {
		CreateScanner(typ reflect.Type, tags reflect.StructTag, paramKeys []string) (scan RequestScanner, err error)
		Describe(op *openapi.Operation, typ reflect.Type) (err error)
	}
)
