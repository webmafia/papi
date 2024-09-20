package fastapi

import (
	"github.com/valyala/fasthttp"
)

type Ctx struct {
	ctx       *fasthttp.RequestCtx
	paramVals []string
}

func (api *API) acquireCtx(c *fasthttp.RequestCtx) (ctx *Ctx) {
	var ok bool

	if ctx, ok = api.ctxPool.Get().(*Ctx); !ok {
		ctx = new(Ctx)
	}

	ctx.ctx = c

	return
}

func (api *API) releaseCtx(ctx *Ctx) {
	ctx.ctx = nil
	ctx.paramVals = ctx.paramVals[:0]
	api.ctxPool.Put(ctx)
}
