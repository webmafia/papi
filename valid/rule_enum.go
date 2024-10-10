package valid

import (
	"reflect"
	"slices"
	"unsafe"

	"github.com/webbmaffian/papi/registry/scanner"
)

func createEnumValidator(offset uintptr, typ reflect.Type, field string, s string) (validator, error) {
	switch kind := typ.Kind(); kind {

	case reflect.Int:
		return validComparableEnum[int](offset, field, s)

	case reflect.Int8:
		return validComparableEnum[int8](offset, field, s)

	case reflect.Int16:
		return validComparableEnum[int16](offset, field, s)

	case reflect.Int32:
		return validComparableEnum[int32](offset, field, s)

	case reflect.Int64:
		return validComparableEnum[int64](offset, field, s)

	case reflect.Uint:
		return validComparableEnum[uint](offset, field, s)

	case reflect.Uint8:
		return validComparableEnum[uint8](offset, field, s)

	case reflect.Uint16:
		return validComparableEnum[uint16](offset, field, s)

	case reflect.Uint32:
		return validComparableEnum[uint32](offset, field, s)

	case reflect.Uint64:
		return validComparableEnum[uint64](offset, field, s)

	case reflect.Float32:
		return validComparableEnum[float32](offset, field, s)

	case reflect.Float64:
		return validComparableEnum[float64](offset, field, s)

	case reflect.Array:
		return validArray(offset, typ, field, s, createEnumValidator)

	case reflect.Slice:
		return validSlice(offset, typ, field, s, createEnumValidator)

	case reflect.Pointer:
		return validPointer(offset, typ, field, s, createEnumValidator)

	case reflect.String:
		return validComparableEnum[string](offset, field, s)

	default:
		return nil, notImplemented("enum", kind)
	}
}

func validComparableEnum[T comparable](offset uintptr, field string, s string) (validator, error) {
	var zero T
	var enum []T

	if err := scanner.ScanString(&enum, s); err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *FieldErrors) {
		if val := *(*T)(unsafe.Add(ptr, offset)); val != zero && !slices.Contains(enum, val) {
			errs.Append(FieldError{
				err:    ErrInvalidEnum,
				field:  field,
				expect: s,
			})
		}
	}, nil
}
