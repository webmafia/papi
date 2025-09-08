package papi

// Route information.
type Route[I any, O any] struct {

	// Mandatory route path. Can contain `{params}`.
	Path string

	// An optional description of the route (longer than the `Summary`).
	Description string

	// Mandatory handler that will be called for the route.
	Handler func(c *RequestCtx, in *I, out *O) error

	// Whether the route is deprecated and discouraged.
	Deprecated bool
}
