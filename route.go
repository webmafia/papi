package papi

import (
	"reflect"
	"runtime"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/errors"
	"github.com/webmafia/papi/internal"
	"github.com/webmafia/papi/internal/route"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
	"github.com/webmafia/papi/valid"
)

// Register a group of routes. Any exported methods with a signature of `func(api *papi.API) error` will be called.
// These methods should call either `papi.GET`, `papi.PUT`, `papi.POST`, or `papi.DELETE`.
func (api *API) RegisterRoutes(types ...any) (err error) {
	for i := range types {
		val := reflect.ValueOf(types[i])
		numMethods := val.NumMethod()

		for i := 0; i < numMethods; i++ {
			cb, ok := val.Method(i).Interface().(func(api *API) error)

			if !ok {
				continue
			}

			if err = cb(api); err != nil {
				return
			}
		}
	}

	return
}

// Add a route with GET method. Input will be validated based on OpenAPI schema rules.
func GET[I, O any](api *API, r Route[I, O]) (err error) {
	var route AdvancedRoute[I, O]

	route.fromRoute(&r)
	route.Method = "GET"

	return addRoute(api, route)
}

// Add a route with PUT method. Input will be validated based on OpenAPI schema rules.
func PUT[I, O any](api *API, r Route[I, O]) (err error) {
	var route AdvancedRoute[I, O]

	route.fromRoute(&r)
	route.Method = "PUT"

	return addRoute(api, route)
}

// Add a route with POST method. Input will be validated based on OpenAPI schema rules.
func POST[I, O any](api *API, r Route[I, O]) (err error) {
	var route AdvancedRoute[I, O]

	route.fromRoute(&r)
	route.Method = "POST"

	return addRoute(api, route)
}

// Add a route with DELETE method. Input will be validated based on OpenAPI schema rules.
func DELETE[I, O any](api *API, r Route[I, O]) (err error) {
	var route AdvancedRoute[I, O]

	route.fromRoute(&r)
	route.Method = "DELETE"

	return addRoute(api, route)
}

func addRoute[I, O any](api *API, r AdvancedRoute[I, O]) (err error) {
	if r.Path == "" {
		return ErrMissingRoutePath
	}

	if r.Handler == nil {
		return ErrMissingRouteHandler
	}

	if r.Method == "" {
		return ErrMissingRouteHandler
	}

	pc := internal.Caller(2)
	validate, err := valid.CreateStructValidator[I]()

	if err != nil {
		return
	}

	return api.router.Add(r.Method, r.Path, func(route *route.Route) (err error) {
		var (
			decodeRequest registry.RequestDecoder
			perm          string
			handler       = *(*registry.Handler)(unsafe.Pointer(&r.Handler))
		)

		if decodeRequest, perm, err = api.reg.CreateRequestDecoder(reflect.TypeFor[I](), route.Params, pc); err != nil {
			return
		}

		if err = addToDocs(api, &r, perm, pc); err != nil {
			return
		}

		if handler, err = api.reg.Handler(reflect.TypeFor[O](), "", route.Params, handler); err != nil {
			return
		}

		route.Handler = func(c *fasthttp.RequestCtx) (err error) {
			var (
				in  I
				out O
			)

			if err = decodeRequest(unsafe.Pointer(&in), c); err != nil {
				return
			}

			var errs errors.Errors

			if !validate(&in, &errs) {
				return errs
			}

			return handler(c, unsafe.Pointer(&in), unsafe.Pointer(&out))
		}

		return
	})
}

func addToDocs[I, O any](api *API, r *AdvancedRoute[I, O], perm string, pc *runtime.Func) (err error) {
	if api.opt.OpenAPI == nil {
		return
	}

	iTyp := reflect.TypeFor[I]()
	oTyp := reflect.TypeFor[O]()

	op := &openapi.Operation{
		Id:          r.OperationId,
		Method:      r.Method,
		Summary:     r.Summary,
		Description: r.Description,
		Tags:        r.Tags,
	}

	if api.opt.SecurityScheme != nil {
		if sec := api.opt.SecurityScheme.OperationSecurityDocs(perm); !sec.IsZero() {
			op.Security = append(op.Security, sec)
		}
	}

	if op.Id == "" {
		title, id := internal.ParseName(internal.CallerNameFromFunc(pc))
		op.Id = id

		if op.Summary == "" {
			op.Summary = title
		}
	}

	if len(op.Tags) == 0 {
		op.Tags = []openapi.Tag{
			openapi.NewTag(internal.CallerTypeFromFunc(pc)),
		}
	}

	if err = api.reg.DescribeOperation(op, iTyp, oTyp); err != nil {
		return
	}

	return internal.AddOperationToDocument(api.opt.OpenAPI, r.Path, op)
}
