package valid

import (
	"reflect"
	"unsafe"

	"github.com/webbmaffian/papi/errors"
	"github.com/webbmaffian/papi/internal/constraints"
	"github.com/webbmaffian/papi/internal/scanner"
)

func createMaxValidator(offset uintptr, typ reflect.Type, field string, s string) (validator, error) {
	switch kind := typ.Kind(); kind {

	case reflect.Int:
		return validNumMax[int](offset, field, s)

	case reflect.Int8:
		return validNumMax[int8](offset, field, s)

	case reflect.Int16:
		return validNumMax[int16](offset, field, s)

	case reflect.Int32:
		return validNumMax[int32](offset, field, s)

	case reflect.Int64:
		return validNumMax[int64](offset, field, s)

	case reflect.Uint:
		return validNumMax[uint](offset, field, s)

	case reflect.Uint8:
		return validNumMax[uint8](offset, field, s)

	case reflect.Uint16:
		return validNumMax[uint16](offset, field, s)

	case reflect.Uint32:
		return validNumMax[uint32](offset, field, s)

	case reflect.Uint64:
		return validNumMax[uint64](offset, field, s)

	case reflect.Float32:
		return validNumMax[float32](offset, field, s)

	case reflect.Float64:
		return validNumMax[float64](offset, field, s)

	case reflect.Array:
		return validArray(offset, typ, field, s, createMaxValidator)

	case reflect.Slice:
		return validSliceMax(offset, field, s)

	case reflect.Pointer:
		return validPointer(offset, typ, field, s, createMaxValidator)

	case reflect.String:
		return validStringMax(offset, field, s)

	default:
		return nil, notImplemented("max", kind)
	}
}

func validNumMax[T constraints.Number](offset uintptr, field string, s string) (validator, error) {
	var max T

	if err := scanner.ScanString(&max, s); err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *errors.Errors) {
		if val := *(*T)(unsafe.Add(ptr, offset)); val != 0 && val > max {
			errs.Append(ErrTooHigh.Explained(field, s))
		}
	}, nil
}

func validSliceMax(offset uintptr, field string, s string) (validator, error) {
	var max int

	if err := scanner.ScanString(&max, s); err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *errors.Errors) {
		if l := sliceLen(unsafe.Add(ptr, offset)); l != 0 && l > max {
			errs.Append(ErrTooLong.Explained(field, s))
		}
	}, nil
}

func validStringMax(offset uintptr, field string, s string) (validator, error) {
	var max int

	if err := scanner.ScanString(&max, s); err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *errors.Errors) {
		if l := stringLen(unsafe.Add(ptr, offset)); l != 0 && l > max {
			errs.Append(ErrTooLong.Explained(field, s))
		}
	}, nil
}
