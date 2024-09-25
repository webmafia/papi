package scanner

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
)

type (
	RequestScanner        func(p unsafe.Pointer, c *fasthttp.RequestCtx) error
	RequestScannerCreator interface {
		CreateScanner(paramKeys []string, tags reflect.StructTag) (scan RequestScanner, err error)
	}
)
