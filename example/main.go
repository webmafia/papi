package main

import (
	"fmt"
	"log"
	"mime/multipart"

	"github.com/webmafia/fastapi"
)

type userRoutes struct{}

func (r userRoutes) GetUser(api *fastapi.API[User]) (err error) {
	type req struct {
		Id    int `param:"id"`
		Limit int `query:"limit"`
	}

	return fastapi.AddRoute(api, fastapi.Route[User, req, User]{
		Method:  "GET",
		Path:    "/users/{id}",
		Summary: "Get user by ID",

		Handler: func(ctx *fastapi.Ctx[User], req *req, resp *User) (err error) {
			resp.ID = req.Id
			resp.Name = "helluuu"

			return
		},
	})
}

func (r userRoutes) ListUsers(api *fastapi.API[User]) (err error) {
	type req struct {
		Status string `query:"status"`
	}

	return fastapi.AddRoute(api, fastapi.Route[User, req, fastapi.List[User]]{
		Method:  "GET",
		Path:    "/users",
		Summary: "List all users",

		Handler: func(ctx *fastapi.Ctx[User], req *req, resp *fastapi.List[User]) (err error) {
			resp.Write(&User{ID: 999, Name: req.Status})
			resp.Write(&User{ID: 998, Name: "Foobaz"})
			resp.Meta.Total = 123

			return
		},
	})
}

func (r userRoutes) CreateUser(api *fastapi.API[User]) (err error) {
	type req struct {
		// Body io.Reader
		Body User
	}

	return fastapi.AddRoute(api, fastapi.Route[User, req, User]{
		Method:  "POST",
		Path:    "/users",
		Summary: "Create user",

		Handler: func(ctx *fastapi.Ctx[User], req *req, resp *User) (err error) {
			// buf, err := io.ReadAll(req.Body)
			// _ = buf
			*resp = req.Body
			resp.ID = 101

			return
		},
	})
}

func (r userRoutes) UploadFile(api *fastapi.API[User]) (err error) {
	type req struct {
		Body *multipart.Form // TODO: Also accept *multipart.File
	}

	return fastapi.AddRoute(api, fastapi.Route[User, req, User]{
		Method:  "POST",
		Path:    "/files",
		Summary: "Upload file",

		Handler: func(ctx *fastapi.Ctx[User], req *req, resp *User) (err error) {
			f := req.Body.File
			fmt.Printf("%#v\n", f)

			return
		},
	})
}

func main() {
	api := fastapi.New[User]()

	if err := api.RegisterRoutes(userRoutes{}); err != nil {
		panic(err)
	}

	log.Println("Listening...")

	if err := api.ListenAndServe("127.0.0.1:3001"); err != nil {
		panic(err)
	}
}
