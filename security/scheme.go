package security

import (
	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/openapi"
)

type Scheme interface {
	SecurityDocs() openapi.SecurityScheme
	OperationSecurityDocs(permTag string) openapi.SecurityRequirement
	// OperationSecurityHandler(typ reflect.Type, perm Permission) (handler func(p unsafe.Pointer, c *fasthttp.RequestCtx) error, modTag string, err error)
	UserRoles(c *fasthttp.RequestCtx) (roles []string, err error)
}
