package fastapi

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"github.com/webmafia/fastapi/pool/json"
	"github.com/webmafia/fastapi/route"
	"github.com/webmafia/fastapi/scanner"
	"github.com/webmafia/fastapi/scanner/request"
	"github.com/webmafia/fastapi/spec"
)

type API struct {
	router   route.Router
	server   fasthttp.Server
	scanners *scanner.Registry
	json     *json.Pool
	docs     *spec.Document
	opt      Options
}

type Options struct {
	JsonAPI jsoniter.API
}

func (opt *Options) setDefaults() {
	if opt.JsonAPI == nil {
		opt.JsonAPI = jsoniter.ConfigFastest
	}
}

func New(opt ...Options) (api *API, err error) {
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
	api.scanners = scanner.NewRegistry(func(r *scanner.Registry) {
		var creator scanner.RequestScannerCreator

		if creator, err = request.NewRequestScanner(r, api.json); err == nil {
			r.RegisterDefaultRequestScanner(creator)
		}
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
	s := api.json.AcquireStream(w)
	defer api.json.ReleaseStream(s)

	// api.docs.JsonEncode(s)

	if err := s.Error; err != nil {
		return err
	}

	return s.Flush()
}
