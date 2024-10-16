package papi

import (
	"github.com/webbmaffian/papi/openapi"
)

// Route information.
type Route[I any, O any] struct {

	// Mandatory route path. Can contain `{params}`.
	Path string

	// An optional description of the route (longer than the `Summary`).
	Description string

	// Optional OpenAPI tags.
	Tags []*openapi.Tag

	// Mandatory handler that will be called for the route.
	Handler func(c *RequestCtx, in *I, out *O) error
}
