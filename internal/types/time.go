package types

import (
	"reflect"
	"time"
	"unsafe"

	"github.com/webbmaffian/papi/internal/registry"
	"github.com/webbmaffian/papi/openapi"
)

func TimeType() registry.TypeRegistrar {
	return timeType{}
}

type timeType struct{}

func (t timeType) Type() reflect.Type {
	return reflect.TypeOf((*time.Time)(nil)).Elem()
}

func (t timeType) TypeDescription(reg *registry.Registry) registry.TypeDescription {
	return registry.TypeDescription{
		Schema: func(_ reflect.StructTag) (openapi.Schema, error) {
			return &openapi.Ref{
				Name: "Timestamp",
				Schema: &openapi.String{
					Format: "RFC3339",
				},
			}, nil
		},
		Decoder: func(_ reflect.StructTag) (registry.Decoder, error) {
			return func(p unsafe.Pointer, s string) (err error) {
				ptr := (*time.Time)(p)
				parsed, err := time.Parse(time.RFC3339, s)

				if err == nil {
					*ptr = parsed
				}

				return
			}, nil
		},
	}
}
