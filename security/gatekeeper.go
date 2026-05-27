package security

import (
	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/openapi"
)

type Gatekeeper interface {

	// Describes the security scheme. This is a general description about how to authenticate
	// in the application.
	SecurityScheme() openapi.SecurityScheme

	// Describes the security requirement per route based on its permission tag (if any).
	SecurityRequirement(perm Permission) openapi.SecurityRequirement

	// Handler that always will be called exactly one (1) time at the beginning of any request,
	// regardless permission or policy. Good for e.g. setting a user value on the context.
	BeforeRequest(c *fasthttp.RequestCtx) error

	// Whether permission tags on routes is optional.
	OptionalPermTag() bool

	// Checks the permission and sets any policy. Any returned error will result in "403 Forbidden".
	CheckPermission(c *fasthttp.RequestCtx, perm Permission, cond Cond) error
}
