package registry

import (
	"reflect"

	"github.com/webmafia/fastapi/registry/value"
)

type CreateValueScanner func(tags reflect.StructTag) (scan value.ValueScanner, err error)
