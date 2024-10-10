package valid

import (
	"reflect"
	"unsafe"
)

type validator func(ptr unsafe.Pointer, errs *FieldErrors)
type validatorCreator func(offset uintptr, typ reflect.Type, field string, s string) (validator, error)
