package types

import (
	"reflect"
	"time"
	"unsafe"

	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
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
		Parser: func(_ reflect.StructTag) (registry.Parser, error) {
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
