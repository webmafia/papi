package scanner

import (
	"reflect"

	"github.com/webmafia/fastapi/scanner/value"
)

type CreateValueScanner func(tags reflect.StructTag) (scan value.ValueScanner, err error)
