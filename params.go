package fastapi

import (
	"github.com/valyala/fasthttp"
	"github.com/webmafia/fastapi/route"
)

const paramsKey = "params"

func RequestParams(c *fasthttp.RequestCtx) *route.Params {
	if params, ok := c.UserValue(paramsKey).(*route.Params); ok {
		return params
	}

	return route.NilParams
}

func setRequestParams(c *fasthttp.RequestCtx, params *route.Params) {
	c.SetUserValue(paramsKey, params)
}
