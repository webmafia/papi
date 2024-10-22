package valid

import (
	"reflect"
	"unsafe"

	"github.com/webmafia/papi/errors"
)

func createRequiredValidator(offset uintptr, typ reflect.Type, field string) (validator, error) {
	isZero, err := createZeroChecker(typ)

	if err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *errors.Errors) {
		if isZero(unsafe.Add(ptr, offset)) {
			errs.Append(ErrRequired.Explained(field, "any value"))
		}
	}, nil
}
