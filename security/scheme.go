package security

import (
	"reflect"
	"runtime"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/openapi"
)

type Scheme interface {
	SecurityDocs() openapi.SecurityScheme
	OperationSecurityDocs(permTag string) openapi.SecurityRequirement
	OperationSecurityHandler(typ reflect.Type, permTag string, caller *runtime.Func) (handler func(p unsafe.Pointer, c *fasthttp.RequestCtx) error, modTag string, err error)
}
