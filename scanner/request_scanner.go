package scanner

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
)

type (
	RequestScanner        func(p unsafe.Pointer, c *fasthttp.RequestCtx) error
	RequestScannerCreator interface {
		CreateScanner(typ reflect.Type, tags reflect.StructTag, paramKeys []string) (scan RequestScanner, err error)
	}
)
