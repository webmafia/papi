package fastapi

import "github.com/webmafia/fastapi/spec"

type Route[T any, I any, O any] struct {
	Method      Method
	Path        string
	Summary     string
	Description string
	Tags        []*spec.Tag
	Handler     func(ctx *Ctx[T], in *I, out *O) error
}
