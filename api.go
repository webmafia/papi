package fastapi

import (
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type API[U any] struct {
	router   Router[U]
	ctxPool  sync.Pool
	jsoniter jsoniter.API
}

func New[U any]() *API[U] {
	return &API[U]{
		jsoniter: jsoniter.ConfigFastest,
	}
}

func (api *API[U]) handler(c *fasthttp.RequestCtx) {
	ctx := api.acquireCtx(c)
	defer api.releaseCtx(ctx)

	cb, params, ok := api.router.Lookup(c.Method(), c.Path(), &ctx.paramVals)

	if !ok {
		// TODO: Proper JSON response
		c.NotFound()
		return
	}

	if len(params) != len(ctx.paramVals) {
		// TODO: Proper JSON response
		c.Response.Reset()
		c.SetStatusCode(fasthttp.StatusInternalServerError)
		c.SetBodyString("params count mismatch")
		return
	}

	if err := cb(ctx); err != nil {
		// TODO: Proper error message
		c.Error(err.Error(), 500)
	}
}

func (api *API[U]) ListenAndServe(addr string) error {
	return fasthttp.ListenAndServe(addr, api.handler)
}
