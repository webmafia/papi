package fastapi

import (
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fastapi/router"
)

type API[U any] struct {
	router  *router.Router
	ctxPool sync.Pool
}

func New[U any]() *API[U] {
	return &API[U]{
		router: router.New(),
	}
}

func (api *API[U]) handler(c *fasthttp.RequestCtx) {
	ctx := api.acquireCtx(c)
	defer api.releaseCtx(ctx)

	ptr := api.router.Lookup(c.Method(), c.Path(), &ctx.params)

	if ptr == nil {
		// TODO: Proper JSON response
		c.NotFound()
		return
	}

	cb := *(*func(*Ctx[U]) error)(ptr)

	if err := cb(ctx); err != nil {
		// TODO: Proper error message
		c.Error(err.Error(), 500)
	}
}

func (api *API[U]) ListenAndServe(addr string) error {
	return fasthttp.ListenAndServe(addr, api.handler)
}
