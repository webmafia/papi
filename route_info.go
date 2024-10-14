package papi

import (
	"github.com/webbmaffian/papi/openapi"
)

type AdvancedRoute[I any, O any] struct {
	OperationId string
	Method      string
	Path        string
	Summary     string
	Description string
	Tags        []*openapi.Tag
	Handler     func(c *RequestCtx, in *I, out *O) error
}

func (adv *AdvancedRoute[I, O]) fromRoute(r *Route[I, O]) {
	adv.Path = r.Path
	adv.Description = r.Description
	adv.Tags = r.Tags
	adv.Handler = r.Handler
}

type Route[I any, O any] struct {
	Path        string
	Description string
	Tags        []*openapi.Tag
	Handler     func(c *RequestCtx, in *I, out *O) error
}
