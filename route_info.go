package fastapi

import (
	"github.com/valyala/fasthttp"
	"github.com/webmafia/fastapi/spec"
)

type Route[I any, O any] struct {
	Method      Method
	Path        string
	Summary     string
	Description string
	Tags        []*spec.Tag
	Handler     func(c *fasthttp.RequestCtx, in *I, out *O) error
}
