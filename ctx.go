package fastapi

import (
	"github.com/valyala/fasthttp"
)

type Ctx[U any] struct {
	ctx       *fasthttp.RequestCtx
	paramVals []string
	User      U
}

func (api *API[U]) acquireCtx(c *fasthttp.RequestCtx) (ctx *Ctx[U]) {
	var ok bool

	if ctx, ok = api.ctxPool.Get().(*Ctx[U]); !ok {
		ctx = new(Ctx[U])
	}

	ctx.ctx = c

	return
}

func (api *API[U]) releaseCtx(ctx *Ctx[U]) {
	ctx.ctx = nil
	ctx.paramVals = ctx.paramVals[:0]
	api.ctxPool.Put(ctx)
}
