# papi
<img alt="Papi" src="./docs/papi.webp" align="right" />

Performant API framework in Go.
- Reduced boilerplate
- Only reflection on startup - no reflection during runtime
- Highly optimized, with almost zero allocations during runtime
- Static typing of routes
- Automatic generation of OpenAPI documentation
- Automatic validation (based on OpenAPI schema rules)
- Encourages dependency injection

**WARNING: This hasn't reached version 1.0.0 and is not production ready yet.**

## Installation
```sh
go get github.com/webmafia/papi
```

## Usage
See the [example](./example) for a full example of how it's used.

### Routing
Papi uses generics in routes to leverage static typing - this also makes it possible to generate an OpenAPI documentation automatically that is always 100% up to date with the code. As generics in Go only can be used on public types and functions, the following methods are public in the package:

```go
papi.GET[I, O any](api *papi.API, r papi.Route[I, O]) (err error)
papi.PUT[I, O any](api *papi.API, r papi.Route[I, O]) (err error)
papi.POST[I, O any](api *papi.API, r papi.Route[I, O]) (err error)
papi.DELETE[I, O any](api *papi.API, r papi.Route[I, O]) (err error)
```
The `I` and `O` generic types are input (request) and output (response), respectively. 

It might look strange at first, but the resulting code gets pretty neat:
```go
type req struct {}

papi.GET(api, papi.Route[req, domain.User]{
	Path: "/users/{id}",

	// A handler always accepts a request- and response type, and returns any error occured.
	Handler: func(ctx *papi.RequestCtx, req *req, resp *domain.User) (err error) {
		resp.ID = 123
		resp.Name = "John Doe"

		return
	},
})
```

By passing pointers of the request and response to the handler, no allocation nor unnecessary copying is needed. The response is often domain model structs, but can be any type.

But how about the request input? In the example above it's an empty struct, but let's explore this in the next section.

### Request input
The request input type can have any name, but it must always be a struct. This allows us to use struct tags for some magic:
```go
type req struct{
	Id int `param:"id"`
}
```

If you look at the previous example, you'll see that the `Path` field contains a parameter in the format `{id}`. As we've tagged our `Id` field above with `param:"id"`, any value passed in the path will end up here. Also, as the type of the field is an `int`, only integers will be accepted - this is validated automatically.

The following tags are supported in the request input:

| Tag           | Meaning                            | Example source  | Example destination      |
| ------------- | ---------------------------------- | --------------- | ------------------------ |
| `param:"*"`   | URL parameters                     | `/users/{id}`   | `123`                    |
| `query:"*"`   | Search query parameters            | `?foo=bar,baz`  | `[]string{"bar","baz"}`  |
| `body:"json"` | PUT and POST bodies in JSON format | `{"foo":"bar"}` | `MyStruct{ Foo: "bar" }` |

Note that string types are not copied, which means that any values in `req` must not be used outside the handler.

### Validation
Sometimes the type is not enought. That's why we support OpenAPI's schema rules. Take this example:
```go
type req struct{
	OrderBy string `query:"orderby" enum:"name,email"`
	Order   string `query:"order" enum:"asc,desc"`
}
```

The following validation tags are supported in the request input (as well as in any nested structs):

| Tag                | Int / Float            | String                 | Slice                              | Array                |
| ------------------ | ---------------------- | ---------------------- | ---------------------------------- | -------------------- |
| `min:"*"`          | Minimum value          | Minimum length         | Minimum length                     | -                    |
| `max:"*"`          | Maximum value          | Maximum length         | Maximum length                     | -                    |
| `enum:"*,*,*"`     | One of specific values | One of specific values | -                                  | -                    |
| `pattern:"*"`      | -                      | Regular expression     | -                                  | -                    |
| `default:"*"`      | Sets default if zero   | Sets default if zero   | Sets default if zero               | Sets default if zero |
| `flags:"required"` | Must be non-zero       | Must be non-zero       | Must have at least 1 non-zero item | Must be non-zero     |

Please note:
- If slices and arrays don't support a tag, it's passed to their children.
- Pointers to any type is only required to be non-nil when required.

### Routing groups & OpenAPI operations
When creating an API you usually want to inject any dependencies, e.g. a User service for any user-related routes - or "operations" as they are called in the OpenAPI specfication. Also, each operation is required to have an API-unique identiier (Operation ID), and is usually grouped by a tag.

Papi solves all this with what we call a routing group, which basically is an arbitrary struct with methods matching the `func(*papi.API) error` signature:

```go
type Users struct{}

func (r Users) GetUserByID(api *papi.API) (err error) {
	type req struct {
		Id int `param:"id"`
	}

	return papi.GET(api, papi.Route[req, User]{
		Path: "/users/{id}",

		Handler: func(ctx *papi.RequestCtx, req *req, resp *domain.User) (err error) {
			resp.ID = 123
			resp.Name = "John Doe"

			return
		},
	})
}

func main() {
	// API initialization and error handling is left out for brevity

	err := api.RegisterRoutes(Users{})
}
```

What happens here:
- As `GetUserByID` matches the `func(*papi.API) error`, this will be called on registration.
- A valid OpenAPI Operation ID will be generated from the method's name, resulting in `get-users-by-id`.
- A descriptive summary of the route will also be generated from the method's name, result in `Get user by ID`.
- The `req` type won't leak outside the route.
- All OpenAPI operations will be assigned a tag matching the group's name, in this case `Users`.
- We are able to inject any dependency into the `Users` struct, and use them in the routes.