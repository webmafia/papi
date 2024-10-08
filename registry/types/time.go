package types

import (
	"reflect"
	"time"
	"unsafe"

	"github.com/webmafia/fastapi/openapi"
	"github.com/webmafia/fastapi/registry/value"
)

type Time struct{}

func (t Time) Type() reflect.Type {
	return reflect.TypeOf((*time.Time)(nil)).Elem()
}

func (t Time) CreateScanner(_ reflect.StructTag) (scan value.ValueScanner, err error) {
	return func(p unsafe.Pointer, s string) (err error) {
		ptr := (*time.Time)(p)
		parsed, err := time.Parse(time.RFC3339, s)

		if err == nil {
			*ptr = parsed
		}

		return
	}, nil
}

func (t Time) Describe(_ reflect.StructTag) openapi.Schema {
	return openapi.Schema{
		Type:   openapi.String,
		Format: "RFC3339",
	}
}
