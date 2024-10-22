package valid

import (
	"reflect"
	"unsafe"

	"github.com/webmafia/papi/errors"
)

// Validates T and returns a boolean whether it's valid
type StructValidator[T any] func(ptr *T, errs *errors.Errors) bool
type structValidator func(ptr unsafe.Pointer, errs *errors.Errors) bool
type validator func(ptr unsafe.Pointer, errs *errors.Errors)
type validatorCreator func(offset uintptr, typ reflect.Type, field string, s string) (validator, error)

type validators []validator

func (valids *validators) append(valid validator) {
	*valids = append(*valids, valid)
}

func (valids *validators) compile() (valid structValidator, err error) {
	*valids = compactSlice(*valids)

	return func(ptr unsafe.Pointer, errs *errors.Errors) bool {
		for _, valid := range *valids {
			valid(ptr, errs)
		}

		return !errs.HasError()
	}, nil
}
