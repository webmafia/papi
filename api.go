package fastapi

import (
	"io"
	"reflect"
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fastapi/internal/jsonpool"
	"github.com/webmafia/fastapi/spec"
	"github.com/webmafia/fastapi/spec/schema"
)

type API[U any] struct {
	router  Router[U]
	ctxPool sync.Pool
	server  fasthttp.Server
	docs    *spec.Document
	opt     Options
}

type Options struct {
	OpenAPI spec.OpenAPI
}

func New[U any](opt ...Options) *API[U] {
	api := &API[U]{
		server: fasthttp.Server{
			StreamRequestBody:            true,
			DisablePreParseMultipartForm: true,
		},
		docs: &spec.Document{
			OpenAPI: "3.0.0",
			Schemas: make(map[reflect.Type]schema.Schema),
		},
	}

	if len(opt) > 0 {
		api.opt = opt[0]
	}

	api.server.Handler = api.handler
	api.docs.Info = api.opt.OpenAPI.Info
	api.docs.Servers = api.opt.OpenAPI.Servers

	return api
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
	return api.server.ListenAndServe(addr)
}

func (api *API[U]) WriteOpenAPI(w io.Writer) error {
	s := jsonpool.AcquireStream(w)
	defer jsonpool.ReleaseStream(s)

	api.docs.JsonEncode(s)

	if err := s.Error; err != nil {
		return err
	}

	return s.Flush()
}
