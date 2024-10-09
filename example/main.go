package main

import (
	"log"
	"os"
	"time"

	"github.com/webbmaffian/papi"
	"github.com/webbmaffian/papi/openapi"
)

var (
	Users = openapi.NewTag("users", "Users")
	Files = openapi.NewTag("files", "Files")
)

type userRoutes struct{}

func (r userRoutes) GetUser(api *papi.API) (err error) {
	type req struct {
		Id int `param:"id"`
	}

	return papi.AddRoute(api, papi.Route[req, User]{
		Method:  "GET",
		Path:    "/users/{id}",
		Summary: "Get user by ID",
		Tags:    []*openapi.Tag{Users},

		Handler: func(ctx *papi.RequestCtx, req *req, resp *User) (err error) {
			resp.ID = req.Id
			resp.Name = "helluuu"

			return
		},
	})
}

func (r userRoutes) ListUsers(api *papi.API) (err error) {
	type req struct {
		Status  string    `query:"status"`
		Before  time.Time `query:"before"`
		Limit   int       `query:"limit" min:"0" max:"500"`
		Decimal float64   `query:"decimal"`
	}

	return papi.AddRoute(api, papi.Route[req, papi.List[User]]{
		Method:  "GET",
		Path:    "/users",
		Summary: "List all users",
		Tags:    []*openapi.Tag{Users},

		Handler: func(ctx *papi.RequestCtx, req *req, resp *papi.List[User]) (err error) {
			resp.Write(&User{ID: 999, Name: req.Status, TimeCreated: req.Before})
			resp.Write(&User{ID: 998, Name: "Foobaz", TimeCreated: req.Before})
			resp.Meta.Total = 123

			return
		},
	})
}

func (r userRoutes) CreateUser(api *papi.API) (err error) {
	type req struct {
		// Body io.Reader
		Body User `body:"json"`
	}

	return papi.AddRoute(api, papi.Route[req, User]{
		Method:  "POST",
		Path:    "/users",
		Summary: "Create user",
		Tags:    []*openapi.Tag{Users},

		Handler: func(ctx *papi.RequestCtx, req *req, resp *User) (err error) {
			// buf, err := io.ReadAll(req.Body)
			// _ = buf
			*resp = req.Body
			resp.ID = 101

			return
		},
	})
}

// func (r userRoutes) UploadFile(api *papi.API[User]) (err error) {
// 	type req struct {
// 		Body *multipart.Form // TODO: Also accept *multipart.File
// 	}

// 	return papi.AddRoute(api, papi.Route[User, req, User]{
// 		Method:  "POST",
// 		Path:    "/files",
// 		Summary: "Upload file",
// 		Tags:    []*spec.Tag{Files},

// 		Handler: func(ctx *papi.Ctx[User], req *req, resp *User) (err error) {
// 			f := req.Body.File
// 			fmt.Printf("%#v\n", f)

// 			return
// 		},
// 	})
// }

func main() {
	api, err := papi.NewAPI(papi.Options{
		OpenAPI: openapi.NewDocument(),
		// OpenAPI: spec.OpenAPI{
		// 	Info: spec.Info{
		// 		Title: "Demo API",
		// 		License: spec.License{
		// 			Name: "MIT",
		// 		},
		// 	},
		// 	Servers: []spec.Server{
		// 		{
		// 			Description: "Local",
		// 			Url:         "http://localhost:3001",
		// 		},
		// 	},
		// },
	})

	if err != nil {
		panic(err)
	}

	if err := api.RegisterRoutes(userRoutes{}); err != nil {
		panic(err)
	}

	if err := dumpSpecToFile(api); err != nil {
		panic(err)
	}

	log.Println("Listening...")

	if err := api.ListenAndServe("127.0.0.1:3001"); err != nil {
		panic(err)
	}
}

func dumpSpecToFile(api *papi.API) (err error) {
	log.Println("Dumping OpenAPI spec to file...")
	f, err := os.Create("openapi.json")

	if err != nil {
		return
	}

	defer f.Close()

	return api.WriteOpenAPI(f)
}
