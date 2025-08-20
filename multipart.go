package papi

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
)

var _ registry.TypeDescriber = (*Multipart[struct{}])(nil)

type Multipart[T any] struct {
	Fields T
}

// TypeDescription implements registry.TypeDescriber.
func (m *Multipart[T]) TypeDescription(reg *registry.Registry) registry.TypeDescription {
	return registry.TypeDescription{
		Schema: func(_ reflect.StructTag) (_ openapi.Schema, err error) {
			schema, err := reg.Schema(reflect.TypeFor[T]())

			if err != nil {
				return
			}

			return &openapi.Custom{
				ContentType: "multipart/form-data",
				Schema:      schema,
			}, nil
		},
	}
}

var _ registry.TypeDescriber = (*MultipartFile)(nil)

type MultipartFile struct{}

// TypeDescription implements registry.TypeDescriber.
func (m MultipartFile) TypeDescription(reg *registry.Registry) registry.TypeDescription {
	return registry.TypeDescription{
		Schema: func(_ reflect.StructTag) (schema openapi.Schema, err error) {
			return &openapi.String{
				Format: "binary",
			}, nil
		},
		Handler: func(handler registry.Handler) (registry.Handler, error) {
			return func(c *fasthttp.RequestCtx, ptr unsafe.Pointer) error {
				return nil
			}, nil
		},
		Decoder: func(tags reflect.StructTag) (registry.Decoder, error) {
			return func(p unsafe.Pointer, s string) error {
				return nil
			}, nil
		},
	}
}
