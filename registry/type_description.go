package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/internal/scanner"
	"github.com/webmafia/papi/openapi"
)

var typeDescriber = reflect.TypeFor[TypeDescriber]()

type TypeDescriber interface {
	TypeDescription(reg *Registry) TypeDescription
}

type TypeRegistrar interface {
	Type() reflect.Type
	TypeDescriber
}

// Parser converts a single string (e.g. from query, path, header) into a value.
// Input-only; no ctx access, no response writes.
type Parser = scanner.Scanner

// Binder populates a value from the request context (body, multipart, streaming).
// Input-only; runs before user handler; may short-circuit.
type Binder func(c *fasthttp.RequestCtx, ptr unsafe.Pointer) error

// Responder shapes/streams the response around user code.
// Output-only; must call next() exactly once (unless short-circuiting).
type Responder func(c *fasthttp.RequestCtx, ptr unsafe.Pointer, next func() error) error

var (
	NoParser = func(reflect.StructTag) (Parser, error) {
		return nil, nil
	}
)

// TypeDescription declares how a type behaves in the API framework.
// A type may implement any subset depending on if it's used for input, output, or both.
type TypeDescription struct {

	// Schema describes the type for OpenAPI (structure, format, content-type).
	Schema func(tags reflect.StructTag) (openapi.Schema, error)

	// Parser provides a string→value parser for this type (query/param/header).
	Parser func(tags reflect.StructTag) (Parser, error)

	// Binder provides a ctx-aware binder for this type (body, multipart, etc.).
	// Preferred over Parser if both are defined.
	Binder func(fieldName string, tags reflect.StructTag) (Binder, error)

	// Responder provides an output wrapper/encoder for this type.
	// Runs after user handler to write/stream the response.
	Responder func() (Responder, error)
}

func (t TypeDescription) IsZero() bool {
	return t.Schema == nil && t.Parser == nil && t.Binder == nil && t.Responder == nil
}
