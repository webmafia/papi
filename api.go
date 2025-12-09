package papi

import (
	"context"
	"io"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
	"github.com/webmafia/papi/errors"
	"github.com/webmafia/papi/internal/iterate"
	"github.com/webmafia/papi/internal/json"
	"github.com/webmafia/papi/internal/route"
	"github.com/webmafia/papi/internal/types"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
	"github.com/webmafia/papi/security"
)

type API struct {
	router route.Router
	server fasthttp.Server
	reg    *registry.Registry
	opt    Options
}

// API options.
type Options struct {

	// An optional (but recommended) OpenAPI document. Provided document must be unused, a.k.a. have no registered operations,
	// and will be filled with documentation for all routes.
	OpenAPI *openapi.Document

	// Any errors occured will be passed through this callback, where it has the chance to transform the error to an
	// `errors.ErrorDocumentor` (if not already). Any error that isn't transformed will be replaced with a general error message.
	TransformError func(err error) errors.ErrorDocumentor

	// Header for Cross-Origin Resource Sharing (CORS).
	CORS string
}

func (opt *Options) setDefaults() {
	if opt.TransformError == nil {
		opt.TransformError = func(err error) errors.ErrorDocumentor {
			if e, ok := err.(errors.ErrorDocumentor); ok {
				return e
			}

			return ErrUnknownError.Explained("", err.Error())
		}
	}
}

// Create a new API service.
func NewAPI(reg *registry.Registry, opt ...Options) (api *API, err error) {
	api = &API{
		server: fasthttp.Server{
			StreamRequestBody:            true,
			DisablePreParseMultipartForm: true,
			Name:                         "papi",
		},
		reg: reg,
	}

	if len(opt) > 0 {
		api.opt = opt[0]

		if api.opt.OpenAPI != nil {
			if api.opt.OpenAPI.NumOperations() != 0 {
				return nil, ErrInvalidOpenAPI
			}

			if g := api.reg.Gatekeeper(); g != nil {
				if err = api.opt.OpenAPI.AddSecurityScheme(g.SecurityScheme()); err != nil {
					return
				}
			}
		}
	}

	api.opt.setDefaults()

	api.reg.RegisterType(
		types.TimeType(),
		multipartFiles{},
	)

	api.server.Handler = api.handler

	return
}

func (api *API) HasPermission(roles []string, perm security.Permission) bool {
	return api.reg.Policies().Has(roles, perm)
}

func (api *API) sendError(c *fasthttp.RequestCtx, err errors.ErrorDocumentor) {
	s := json.AcquireStream(c)
	defer json.ReleaseStream(s)

	c.Response.Reset()
	c.SetStatusCode(err.Status())
	api.cors(c)

	err.ErrorDocument(s)
	s.Flush()
}

func (api *API) handler(c *fasthttp.RequestCtx) {
	if api.cors(c) {
		return
	}

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

// Listen on the provided address (e.g. `localhost:3000`).
func (api *API) ListenAndServe(addr string) error {
	return api.server.ListenAndServe(addr)
}

// Close API for new requests, and close all current requests after specified grace period (default 3 seconds).
func (api *API) Close(grace ...time.Duration) error {
	var wait time.Duration

	if len(grace) > 0 {
		wait = grace[0]
	} else {
		wait = 3 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	return api.server.ShutdownWithContext(ctx)
}

// Write API documentation to an `io.Writer`.`
func (api *API) WriteOpenAPI(w io.Writer) error {
	if api.opt.OpenAPI == nil {
		return ErrMissingOpenAPI
	}

	s := jsoniter.Config{
		IndentionStep:                 4,
		MarshalFloatWith6Digits:       true,
		ObjectFieldMustBeSimpleString: true,
	}.Froze().BorrowStream(w)

	if err := api.opt.OpenAPI.JsonEncode(s); err != nil {
		return err
	}

	if err := s.Error; err != nil {
		return err
	}

	return s.Flush()
}

// Register a custom type, that will override any defaults.
func (api *API) RegisterType(typs ...registry.TypeRegistrar) (err error) {
	return api.reg.RegisterType(typs...)
}

// TODO: Real implementation
func (api *API) cors(c *fasthttp.RequestCtx) (abort bool) {
	cors := api.opt.CORS

	if cors == "" {
		return
	}

	if api.opt.CORS == "*" {
		if origin := c.Request.Header.Peek("Origin"); len(origin) > 0 {
			cors = fast.BytesToString(origin)
		}
	} else if strings.IndexByte(api.opt.CORS, ',') >= 0 {
		if origin := fast.BytesToString(c.Request.Header.Peek("Origin")); len(origin) > 0 {
			for _, chunk := range iterate.IterateChunks(api.opt.CORS, ',') {
				if chunk == origin {
					cors = chunk
				}
			}
		}
	}

	c.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	c.Response.Header.Set("Access-Control-Allow-Headers", "authorization,content-type")
	c.Response.Header.Set("Access-Control-Allow-Methods", "HEAD,GET,POST,PUT,DELETE,OPTIONS")
	c.Response.Header.Set("Access-Control-Allow-Origin", cors)

	return c.IsOptions()
}
