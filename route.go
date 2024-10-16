package papi

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webbmaffian/papi/internal"
	"github.com/webbmaffian/papi/openapi"
	"github.com/webbmaffian/papi/registry"
	"github.com/webbmaffian/papi/route"
	"github.com/webbmaffian/papi/valid"
)

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

func GET[I, O any](api *API, r Route[I, O]) (err error) {
	var route AdvancedRoute[I, O]

	route.fromRoute(&r)
	route.Method = "GET"

	if err = addToDocs(api, &route); err != nil {
		return
	}

	return addRoute(api, route)
}

func PUT[I, O any](api *API, r Route[I, O]) (err error) {
	var route AdvancedRoute[I, O]

	route.fromRoute(&r)
	route.Method = "PUT"

	if err = addToDocs(api, &route); err != nil {
		return
	}

	return addRoute(api, route)
}

func POST[I, O any](api *API, r Route[I, O]) (err error) {
	var route AdvancedRoute[I, O]

	route.fromRoute(&r)
	route.Method = "POST"

	if err = addToDocs(api, &route); err != nil {
		return
	}

	return addRoute(api, route)
}

func DELETE[I, O any](api *API, r Route[I, O]) (err error) {
	var route AdvancedRoute[I, O]

	route.fromRoute(&r)
	route.Method = "DELETE"

	if err = addToDocs(api, &route); err != nil {
		return
	}

	return addRoute(api, route)
}

func AddRoute[I, O any](api *API, r AdvancedRoute[I, O]) (err error) {
	if err = addToDocs(api, &r); err != nil {
		return
	}

	return addRoute(api, r)
}

func addRoute[I, O any](api *API, r AdvancedRoute[I, O]) (err error) {
	validate, err := valid.CreateStructValidator[I]()

	if err != nil {
		return
	}

	return api.router.Add(r.Method, r.Path, func(route *route.Route) (err error) {
		var (
			decodeRequest registry.RequestDecoder
			handler       = *(*registry.Handler)(unsafe.Pointer(&r.Handler))
		)

		if decodeRequest, err = api.reg.CreateRequestDecoder(reflect.TypeFor[I](), route.Params); err != nil {
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

			// TODO: Reuse errors from pool, and return any errors to client
			var errs valid.FieldErrors
			validate(&in, &errs)

			if errs.HasError() {
				return errs
			}

			return handler(c, unsafe.Pointer(&in), unsafe.Pointer(&out))
		}

		return
	})
}

func addToDocs[I, O any](api *API, r *AdvancedRoute[I, O]) (err error) {
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

	if op.Id == "" {
		title, id := internal.ParseName(internal.CallerName(2))
		op.Id = id

		if op.Summary == "" {
			op.Summary = title
		}
	}

	if err = api.reg.DescribeOperation(op, iTyp, oTyp); err != nil {
		return
	}

	return api.opt.OpenAPI.AddOperation(r.Path, op)
}
