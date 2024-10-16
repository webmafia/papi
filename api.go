package papi

import (
	"errors"
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	papierr "github.com/webbmaffian/papi/errors"
	"github.com/webbmaffian/papi/openapi"
	"github.com/webbmaffian/papi/pool/json"
	"github.com/webbmaffian/papi/registry"
	"github.com/webbmaffian/papi/registry/types"
	"github.com/webbmaffian/papi/route"
)

type API struct {
	router route.Router
	server fasthttp.Server
	reg    *registry.Registry
	json   *json.Pool
	opt    Options
}

type Options struct {
	JsonAPI        jsoniter.API
	OpenAPI        *openapi.Document
	TransformError func(err error) papierr.ErrorDocumentor
}

func (opt *Options) setDefaults() {
	if opt.JsonAPI == nil {
		opt.JsonAPI = jsoniter.ConfigFastest
	}

	if opt.TransformError == nil {
		opt.TransformError = func(err error) papierr.ErrorDocumentor {
			if e, ok := err.(papierr.ErrorDocumentor); ok {
				return e
			}

			return ErrUnknownError
		}
	}
}

func NewAPI(opt ...Options) (api *API, err error) {
	api = &API{
		server: fasthttp.Server{
			StreamRequestBody:            true,
			DisablePreParseMultipartForm: true,
		},
	}

	if len(opt) > 0 {
		api.opt = opt[0]

		if api.opt.OpenAPI.NumOperations() != 0 {
			return nil, errors.New("there must not be any existing operations in OpenAPI documentation")
		}
	}

	api.opt.setDefaults()
	api.json = json.NewPool(api.opt.JsonAPI)

	if api.reg, err = registry.NewRegistry(api.json); err != nil {
		return
	}

	api.reg.RegisterType(
		types.TimeType(),
	)

	api.server.Handler = api.handler

	return
}

func (api *API) sendError(c *fasthttp.RequestCtx, err papierr.ErrorDocumentor) {
	s := api.json.AcquireStream(c)
	defer api.json.ReleaseStream(s)

	c.Response.Reset()
	c.SetStatusCode(err.Status())

	err.ErrorDocument(s)
	s.Flush()
}

func (api *API) handler(c *fasthttp.RequestCtx) {
	cb, params, ok := api.router.Lookup(c.Method(), c.Path())
	defer api.router.ReleaseParams(params)

	if !ok {
		api.sendError(c, ErrNotFound)
		return
	}

	if !params.Valid() {
		api.sendError(c, ErrInvalidParams)
		return
	}

	route.SetRequestParams(c, params)

	if err := cb(c); err != nil {
		api.sendError(c, api.opt.TransformError(err))
		return
	}
}

func (api *API) ListenAndServe(addr string) error {
	return api.server.ListenAndServe(addr)
}

func (api *API) WriteOpenAPI(w io.Writer) error {
	if api.opt.OpenAPI == nil {
		return errors.New("no OpenAPI documentation initialized")
	}

	s := api.json.AcquireStream(w)
	defer api.json.ReleaseStream(s)

	if err := api.opt.OpenAPI.JsonEncode(s); err != nil {
		return err
	}

	if err := s.Error; err != nil {
		return err
	}

	return s.Flush()
}
