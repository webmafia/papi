package papi

import (
	"bytes"
	"slices"
	"sync"
)

type docsRoute struct{}

func (docsRoute) OpenAPI(api *API) error {
	type req struct {
		Policy struct{} `perm:"-"`
	}

	var mu sync.Mutex
	var json []byte

	return AddRoute(Advanced(api), AdvancedRoute[struct{}, struct{}]{
		Method:         "GET",
		Path:           "/openapi.json",
		HiddenFromDocs: true,
		Handler: func(c *RequestCtx, _, _ *struct{}) (err error) {
			mu.Lock()
			defer mu.Unlock()

			if json == nil {
				var buf bytes.Buffer

				if err = api.WriteOpenAPI(&buf); err != nil {
					return
				}

				json = slices.Clip(buf.Bytes())
			}

			c.SetBody(json)

			return
		},
	})
}
