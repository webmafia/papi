package fastapi

type Route[T any, I any, O any] struct {
	Method      Method
	Path        string
	Summary     string
	Description string
	Handler     func(ctx *Ctx[T], in *I, out *O) error
}
