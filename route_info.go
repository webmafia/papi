package fastapi

import (
	"github.com/valyala/fasthttp"
	"github.com/webbmaffian/papi/openapi"
)

type Route[I any, O any] struct {
	Method      Method
	Path        string
	Summary     string
	Description string
	Tags        []*openapi.Tag
	Handler     func(c *fasthttp.RequestCtx, in *I, out *O) error
}
