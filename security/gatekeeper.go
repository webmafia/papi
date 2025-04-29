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

	// Handler that always will be called exactly one (1) time at the beginning of any request,
	// regardless permission or policy. Good for e.g. setting a user value on the context.
	PreRequest(c *fasthttp.RequestCtx) error

	// Returns the roles that a particular HTTP request has. Will only be called on routes with
	// a permission requirement set.
	UserRoles(c *fasthttp.RequestCtx) (roles []string, err error)

	// Whether permission tags on routes is optional.
	OptionalPermTag() bool
}
