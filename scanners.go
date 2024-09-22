package fastapi

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fastapi/internal"
)

func RegisterRequestScanner[T any](api *API, fn func(v *T, reqCtx *fasthttp.RequestCtx, paramVals []string) error) {
	typ := internal.ReflectType[T]()

	if fn == nil {
		api.scanners.set(typ, nil)
		return
	}

	api.scanners.set(typ, func(p unsafe.Pointer, reqCtx *fasthttp.RequestCtx, paramVals []string) error {
		return fn((*T)(p), reqCtx, paramVals)
	})
}

type scanners struct {
	registry sync.Map
}

func (sc *scanners) get(typ reflect.Type) (scan RequestScanner, ok bool) {
	v, ok := sc.registry.Load(typ)

	if ok {
		scan, ok = v.(RequestScanner)
	}

	return
}

func (sc *scanners) set(typ reflect.Type, scan RequestScanner) {
	if scan == nil {
		sc.registry.Delete(typ)
	} else {
		sc.registry.Store(typ, scan)
	}
}