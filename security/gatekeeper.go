package security

import (
	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/openapi"
)

type Gatekeeper interface {

	// OpenAPI description of the security scheme.
	SecurityScheme() openapi.SecurityScheme

	// OpenAPI description of the security requirement.
	SecurityRequirement(perm Permission) openapi.SecurityRequirement

	// Returns the roles that a particular HTTP request has.
	UserRoles(c *fasthttp.RequestCtx) (roles []string, err error)

	// Whether permission tags on routes is optional.
	OptionalPermTag() bool
}
