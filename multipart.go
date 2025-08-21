package papi

import (
	"log"
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
)

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
		Parser: registry.NoParser,
		Binder: func(tags reflect.StructTag) (registry.Binder, error) {
			return func(c *fasthttp.RequestCtx, ptr unsafe.Pointer) error {
				log.Println("tadaaa")
				return nil
			}, nil
		},
		// Handler: func(handler registry.Handler) (registry.Handler, error) {
		// 	return func(c *fasthttp.RequestCtx, ptr unsafe.Pointer) error {
		// 		return nil
		// 	}, nil
		// },
		// Decoder: func(tags reflect.StructTag) (registry.Decoder, error) {
		// 	return func(p unsafe.Pointer, s string) error {
		// 		return nil
		// 	}, nil
		// },
	}
}
