package valid

import (
	"reflect"
	"unsafe"

	"github.com/webbmaffian/papi/internal/constraints"
	"github.com/webbmaffian/papi/registry/scanner"
)

func appendMinValidators(offset uintptr, typ reflect.Type, field string, s string) (validator, error) {
	switch kind := typ.Kind(); kind {

	case reflect.Int:
		return validNumMin[int](offset, field, s)

	case reflect.Int8:
		return validNumMin[int8](offset, field, s)

	case reflect.Int16:
		return validNumMin[int16](offset, field, s)

	case reflect.Int32:
		return validNumMin[int32](offset, field, s)

	case reflect.Int64:
		return validNumMin[int64](offset, field, s)

	case reflect.Uint:
		return validNumMin[uint](offset, field, s)

	case reflect.Uint8:
		return validNumMin[uint8](offset, field, s)

	case reflect.Uint16:
		return validNumMin[uint16](offset, field, s)

	case reflect.Uint32:
		return validNumMin[uint32](offset, field, s)

	case reflect.Uint64:
		return validNumMin[uint64](offset, field, s)

	case reflect.Float32:
		return validNumMin[float32](offset, field, s)

	case reflect.Float64:
		return validNumMin[float64](offset, field, s)

	case reflect.Slice:
		return validSliceMin(offset, field, s)

	case reflect.String:
		return validStringMin(offset, field, s)

	default:
		return nil, notImplemented("min", kind)
	}
}

func validNumMin[T constraints.Number](offset uintptr, field string, s string) (validator, error) {
	var min T

	if err := scanner.ScanString(&min, s); err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *FieldErrors) {
		if val := *(*T)(unsafe.Add(ptr, offset)); val != 0 && val < min {
			errs.Append(FieldError{
				err:    ErrBelowMin,
				field:  field,
				expect: s,
			})
		}
	}, nil
}

func validSliceMin(offset uintptr, field string, s string) (validator, error) {
	var min int

	if err := scanner.ScanString(&min, s); err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *FieldErrors) {
		if l := sliceLen(unsafe.Add(ptr, offset)); l != 0 && l < min {
			errs.Append(FieldError{
				err:    ErrBelowMin,
				field:  field,
				expect: s,
			})
		}
	}, nil
}

func validStringMin(offset uintptr, field string, s string) (validator, error) {
	var min int

	if err := scanner.ScanString(&min, s); err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *FieldErrors) {
		if l := stringLen(unsafe.Add(ptr, offset)); l != 0 && l < min {
			errs.Append(FieldError{
				err:    ErrBelowMin,
				field:  field,
				expect: s,
			})
		}
	}, nil
}
