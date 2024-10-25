package valid

import (
	"reflect"
	"unsafe"

	"github.com/webmafia/papi/errors"
	"github.com/webmafia/papi/internal/scanner"
)

func createDefaultValidator(offset uintptr, typ reflect.Type, field string, s string) (validator, error) {
	isZero, err := createZeroChecker(typ)

	if err != nil {
		return nil, err
	}

	scan, err := scanner.CreateScanner(typ)

	if err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *errors.Errors) {
		ptr = unsafe.Add(ptr, offset)

		if isZero(ptr) {
			if err := scan(ptr, s); err != nil {
				errs.Append(ErrFailedDefault.Detailed(err.Error(), field))
			}
		}
	}, nil
}
