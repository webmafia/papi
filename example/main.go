package main

import "github.com/webmafia/fastapi"

type userRoutes struct{}

func (r userRoutes) GetUser(api *fastapi.API[User]) {
	type req struct {
		Id string `param:"id"`
	}

	type resp struct {
		Body User
	}

	fastapi.AddRoute(api, fastapi.Route[User, req, resp]{
		Method:  "GET",
		Path:    "/users/{id}",
		Summary: "Get user by ID",

		Handler: func(ctx *fastapi.Ctx[User], req *req, resp *resp) (err error) {
			// Do something

			return
		},
	})
}

func (r userRoutes) ListUsers(api *fastapi.API[User]) {
	type req struct {
		Status string `query:"status"`
	}

	fastapi.AddRoute(api, fastapi.Route[User, req, fastapi.Stream[User]]{
		Method:  "GET",
		Path:    "/users",
		Summary: "List all users",

		Handler: func(ctx *fastapi.Ctx[User], req *req, resp *fastapi.Stream[User]) (err error) {
			// Do something

			return
		},
	})
}
