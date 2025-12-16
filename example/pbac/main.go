package main

import (
	"log"
	"os"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi"
	"github.com/webmafia/papi/example"
	"github.com/webmafia/papi/internal"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
	"github.com/webmafia/papi/security"
)

type User = example.User
type Users struct{}

type Policy struct {
	Foo int
	Bar int
}

func (r Users) GetUserByID(api *papi.API) (err error) {
	type req struct {
		Policy Policy `perm:"foobar"`
		Id     int    `param:"id"`
	}

	return papi.GET(api, papi.Route[req, User]{
		Path: "/users/{id}",

		Handler: func(ctx *papi.RequestCtx, req *req, resp *User) (err error) {
			resp.ID = req.Id
			resp.Name = "John Doe"

			log.Println(req.Policy.Foo, req.Policy.Bar)

			return
		},
	})
}

var _ security.RouteGatekeeper = (*gatekeeper)(nil)

type gatekeeper struct{}

// CheckPermission implements security.RouteGatekeeper.
func (g *gatekeeper) CheckPermission(c *fasthttp.RequestCtx, perm security.Permission, policy internal.Setter) error {
	return policy.Set(&Policy{Foo: 123, Bar: 456})
}

// OptionalPermTag implements security.Gatekeeper.
func (g *gatekeeper) OptionalPermTag() bool {
	return true
}

// PreRequest implements security.Gatekeeper.
func (g *gatekeeper) PreRequest(c *fasthttp.RequestCtx) error {
	return nil
}

// SecurityRequirement implements security.Gatekeeper.
func (g *gatekeeper) SecurityRequirement(perm security.Permission) openapi.SecurityRequirement {
	sec := openapi.SecurityRequirement{
		Name: "token",
	}

	if !perm.IsZero() {
		sec.Scopes = []string{perm.String()}
	}

	return sec
}

// SecurityScheme implements security.Gatekeeper.
func (g *gatekeeper) SecurityScheme() openapi.SecurityScheme {
	return openapi.SecurityScheme{}
}

func main() {
	host := "localhost:3001"
	api, err := papi.NewAPI(registry.NewRegistry(&gatekeeper{}), papi.Options{
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
