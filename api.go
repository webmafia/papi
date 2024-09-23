package fastapi

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"github.com/webmafia/fastapi/internal"
	"github.com/webmafia/fastapi/pool/json"
	"github.com/webmafia/fastapi/route"
	"github.com/webmafia/fastapi/scanner/strings"
	"github.com/webmafia/fastapi/scanner/structs"
	"github.com/webmafia/fastapi/spec"
)

type API struct {
	router   route.Router
	server   fasthttp.Server
	scanners scanners
	// docs    *spec.Document
	opt           Options
	scanInputTags structs.TagScanner
}

type Options struct {
	OpenAPI    spec.OpenAPI
	StringScan *strings.Factory
	JsonPool   *json.Pool
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

	if api.opt.StringScan == nil {
		api.opt.StringScan = strings.NewFactory()
	}

	if api.opt.JsonPool == nil {
		api.opt.JsonPool = json.NewPool(jsoniter.ConfigFastest)
	}

	if api.scanInputTags, err = structs.CreateTagScanner(api.opt.StringScan, internal.ReflectType[inputTags]()); err != nil {
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

	setRequestParams(c, params)

	if err := cb(c); err != nil {
		// TODO: Proper error message
		c.Error(err.Error(), 500)
	}
}

func (api *API) ListenAndServe(addr string) error {
	return api.server.ListenAndServe(addr)
}

func (api *API) WriteOpenAPI(w io.Writer) error {
	s := api.opt.JsonPool.AcquireStream(w)
	defer api.opt.JsonPool.ReleaseStream(s)

	// api.docs.JsonEncode(s)

	if err := s.Error; err != nil {
		return err
	}

	return s.Flush()
}
