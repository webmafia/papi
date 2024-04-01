package fastapi

import "github.com/webmafia/fastapi/router"

type API[U any] struct {
	router *router.Router
}

func New[U any]() *API[U] {
	return &API[U]{
		router: router.New(),
	}
}
