package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/webmafia/papi"
	"github.com/webmafia/papi/openapi"
)

type Users struct{}

func (r Users) GetUserByID(api *papi.API) (err error) {
	type req struct {
		Id int `param:"id"`
	}

	return papi.GET(api, papi.Route[req, User]{
		Path: "/users/{id}",

		Handler: func(ctx *papi.RequestCtx, req *req, resp *User) (err error) {
			resp.ID = req.Id
			resp.Name = "John Doe"

			return
		},
	})
}

func (r Users) ListUsers(api *papi.API) (err error) {
	type req struct {
		Status  string    `query:"status"`
		Before  time.Time `query:"before"`
		Limit   int       `query:"limit" min:"0" max:"500"`
		Decimal float64   `query:"decimal"`
	}

	return papi.GET(api, papi.Route[req, papi.List[User]]{
		Path: "/users",

		Handler: func(ctx *papi.RequestCtx, req *req, resp *papi.List[User]) (err error) {
			resp.Write(&User{ID: 999, Name: req.Status, TimeCreated: req.Before})
			resp.Write(&User{ID: 998, Name: "Foobaz", TimeCreated: req.Before})
			resp.SetTotal(123)

			return
		},
	})
}

func (r Users) CreateUser(api *papi.API) (err error) {
	type req struct {
		Body User `body:"json"`
	}

	return papi.POST(api, papi.Route[req, User]{
		Path: "/users",

		Handler: func(ctx *papi.RequestCtx, req *req, resp *User) (err error) {
			*resp = req.Body
			resp.ID = 101

			return
		},
	})
}

func (r Users) DownloadFile(api *papi.API) (err error) {
	type req struct{}

	return papi.GET(api, papi.Route[req, papi.File[PDF]]{
		Path: "/file",

		Handler: func(ctx *papi.RequestCtx, req *req, resp *papi.File[PDF]) (err error) {
			resp.SetFilename("foobar.pdf")

			// This is obviously invalid JSON, but proves the point.
			_, err = fmt.Fprintf(resp.Writer(), "hello %d", 123)
			return
		},
	})
}

func (r Users) RawJson(api *papi.API) (err error) {
	type req struct {
		Body papi.RawJSON `body:"json"`
	}

	// In this case we're demonstrating that RawJSON can be used for both request and response.
	return papi.POST(api, papi.Route[req, papi.RawJSON]{
		Path: "/raw-json",

		Handler: func(ctx *papi.RequestCtx, req *req, resp *papi.RawJSON) (err error) {

			// Here we just send back the request's JSON body. Please don't do this.
			*resp = req.Body
			return
		},
	})
}

var _ papi.FileType = PDF{}

type PDF struct{}

// Binary implements papi.FileType.
func (PDF) Binary() bool { return true }

// ContentType implements papi.FileType.
func (p PDF) ContentType() string { return "application/pdf" }

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

	if err := api.RegisterRoutes(Users{}); err != nil {
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
