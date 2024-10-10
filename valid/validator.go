package valid

import (
	"reflect"
	"unsafe"
)

type Validator[T any] func(ptr *T, errs *FieldErrors)
type validator func(ptr unsafe.Pointer, errs *FieldErrors)
type validatorCreator func(offset uintptr, typ reflect.Type, field string, s string) (validator, error)

type validators []validator

func (valids *validators) append(valid validator) {
	*valids = append(*valids, valid)
}

func (valids *validators) compile() (valid validator, err error) {
	*valids = compactSlice(*valids)

	return func(ptr unsafe.Pointer, errs *FieldErrors) {
		for _, valid := range *valids {
			valid(ptr, errs)
		}
	}, nil
}
