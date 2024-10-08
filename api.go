package fastapi

import (
	"errors"
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"github.com/webmafia/fastapi/openapi"
	"github.com/webmafia/fastapi/pool/json"
	"github.com/webmafia/fastapi/registry"
	"github.com/webmafia/fastapi/registry/request"
	"github.com/webmafia/fastapi/route"
)

type API struct {
	router route.Router
	server fasthttp.Server
	reg    *registry.Registry
	json   *json.Pool
	opt    Options
}

type Options struct {
	JsonAPI jsoniter.API
	OpenAPI *openapi.Document
}

func (opt *Options) setDefaults() {
	if opt.JsonAPI == nil {
		opt.JsonAPI = jsoniter.ConfigFastest
	}
}

func NewAPI(opt ...Options) (api *API, err error) {
	api = &API{
		server: fasthttp.Server{
			StreamRequestBody:            true,
			DisablePreParseMultipartForm: true,
		},
		// docs: &spec.Document{
		// 	OpenAPI: "3.0.0",
		// 	Schemas: make(map[reflect.Type]schema.Schema),
		// },
	}

	if len(opt) > 0 {
		api.opt = opt[0]
	}

	api.opt.setDefaults()

	api.json = json.NewPool(api.opt.JsonAPI)
	api.reg = registry.NewRegistry(func(r *registry.Registry) (creator registry.RequestScannerCreator) {
		r.RegisterCommonTypes()
		creator, err = request.NewRequestScanner(r, api.json)
		return
	})

	if err != nil {
		return
	}

	api.server.Handler = api.handler
	// api.docs.Info = api.opt.OpenAPI.Info
	// api.docs.Servers = api.opt.OpenAPI.Servers

	return
}

func (api *API) handler(c *fasthttp.RequestCtx) {
	cb, params, ok := api.router.Lookup(c.Method(), c.Path())
	defer api.router.ReleaseParams(params)

	if !ok {
		// TODO: Proper JSON response
		c.NotFound()
		return
	}

	if !params.Valid() {
		// TODO: Proper JSON response
		c.Response.Reset()
		c.SetStatusCode(fasthttp.StatusInternalServerError)
		c.SetBodyString("params count mismatch")
		return
	}

	request.SetRequestParams(c, params)

	if err := cb(c); err != nil {
		// TODO: Proper error message
		c.Error(err.Error(), 500)
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
