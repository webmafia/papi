# papi
<img alt="Papi" src="./docs/papi-mug.webp" align="right" />

Performant API framework in Go.
- Reduced boilerplate
- Only reflection on startup - no reflection during runtime
- Highly optimized, with almost zero allocations during runtime
- Static typing of routes
- Automatic generation of OpenAPI documentation
- Automatic validation (based on OpenAPI schema rules)
- Prepared to be used in hexagonal architectures

## Installation
```sh
go get github.com/webbmaffian/papi
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
	Tags: []*openapi.Tag{Users},

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

| Tag   | Meaning                            | Example source  | Example usage |
| ----- | ---------------------------------- | --------------- | ------------- |
| param | URL parameters                     | `/users/{id}`   | `param:"id"`  |
| query | Search query parameters            | `?foo=bar`      | `query:"foo"` |
| body  | PUT and POST bodies in JSON format | `{"foo":"bar"}` | `body:"json"` |

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

To be continued...