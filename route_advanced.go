package papi

import "github.com/webbmaffian/papi/openapi"

type AdvancedRoute[I any, O any] struct {
	OperationId string
	Method      string
	Path        string
	Summary     string
	Description string
	Tags        []*openapi.Tag
	Handler     func(c *RequestCtx, in *I, out *O) error
}

func (adv *AdvancedRoute[I, O]) fromRoute(r *Route[I, O]) {
	adv.Path = r.Path
	adv.Description = r.Description
	adv.Tags = r.Tags
	adv.Handler = r.Handler
}

type advancedApi struct {
	api *API
}

// Expose an advanced, discouraged API.
func Advanced(api *API) advancedApi {
	return advancedApi{
		api: api,
	}
}

// Add an advanced route. This is more low-level without any real benefits, thus discouraged.
func AddRoute[I, O any](api advancedApi, r AdvancedRoute[I, O]) (err error) {
	return addRoute(api.api, r)
}
