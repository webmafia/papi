package main

import (
	"log"
	"os"
	"time"

	"github.com/webbmaffian/papi"
	"github.com/webbmaffian/papi/openapi"
	"github.com/webbmaffian/papi/registry/types"
)

var (
	Users = openapi.NewTag("users", "Users")
	Files = openapi.NewTag("files", "Files")
)

type userRoutes struct{}

func (r userRoutes) GetUserByID(api *papi.API) (err error) {
	type req struct {
		Id int `param:"id"`
	}

	return papi.GET(api, papi.Route[req, User]{
		Path: "/users/{id}",
		Tags: []*openapi.Tag{Users},

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

	return papi.GET(api, papi.Route[req, types.List[User]]{
		Path: "/users",
		Tags: []*openapi.Tag{Users},

		Handler: func(ctx *papi.RequestCtx, req *req, resp *types.List[User]) (err error) {
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

	return papi.POST(api, papi.Route[req, User]{
		Path: "/users",
		Tags: []*openapi.Tag{Users},

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
	host := "localhost:3001"
	api, err := papi.NewAPI(papi.Options{
		OpenAPI: openapi.NewDocument(
			openapi.Info{
				Title: "Demo API",
				License: openapi.License{
					Name: "MIT",
				},
			},
			openapi.Server{
				Description: "Local",
				Url:         "http://" + host,
			},
		),
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

	log.Println("Listening at", host, "...")

	if err := api.ListenAndServe(host); err != nil {
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
