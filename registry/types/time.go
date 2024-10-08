package types

import (
	"reflect"
	"time"
	"unsafe"

	"github.com/webmafia/fastapi/openapi"
	"github.com/webmafia/fastapi/registry/value"
)

func Time() Type {
	return timeType{}
}

type timeType struct{}

func (t timeType) Type() reflect.Type {
	return reflect.TypeOf((*time.Time)(nil)).Elem()
}

func (t timeType) CreateScanner(_ reflect.StructTag) (scan value.ValueScanner, err error) {
	return func(p unsafe.Pointer, s string) (err error) {
		ptr := (*time.Time)(p)
		parsed, err := time.Parse(time.RFC3339, s)

		if err == nil {
			*ptr = parsed
		}

		return
	}, nil
}

func (t timeType) Describe(_ reflect.StructTag) *openapi.Schema {
	return &openapi.Schema{
		Type:   openapi.String,
		Format: "RFC3339",
	}
}