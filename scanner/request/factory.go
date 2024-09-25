package request

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fastapi/scanner/strings"
)

type Factory struct {
	types sync.Map
}

type (
	RequestScanner       func(p unsafe.Pointer, c *fasthttp.RequestCtx) error
	CreateRequestScanner func(f *strings.Factory, typ reflect.Type, params []string, tags reflect.StructTag) (scan RequestScanner, err error)
)

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) CreateScanner(sf *strings.Factory, typ reflect.Type, params []string, tags reflect.StructTag) (scan RequestScanner, err error) {
	if val, ok := f.types.Load(typ); ok {
		if create, ok := val.(CreateRequestScanner); ok {
			return create(sf, typ, params, tags)
		}
	}

}
