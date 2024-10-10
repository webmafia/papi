package valid

import (
	"reflect"
	"unsafe"
)

func appendRequiredValidators(offset uintptr, typ reflect.Type, field string) (validator, error) {
	isZero, err := createZeroChecker(typ)

	if err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *FieldErrors) {
		if isZero(unsafe.Add(ptr, offset)) {
			errs.Append(FieldError{
				err:    ErrRequired,
				field:  field,
				expect: "any value",
			})
		}
	}, nil
}
