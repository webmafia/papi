package valid

import (
	"reflect"
	"unsafe"
)

func appendRequiredValidators(offset uintptr, typ reflect.Type, field string) (validator, error) {
	switch kind := typ.Kind(); kind {

	case reflect.Int:
		return validComparableRequired[int](offset, field)

	case reflect.Int8:
		return validComparableRequired[int8](offset, field)

	case reflect.Int16:
		return validComparableRequired[int16](offset, field)

	case reflect.Int32:
		return validComparableRequired[int32](offset, field)

	case reflect.Int64:
		return validComparableRequired[int64](offset, field)

	case reflect.Uint:
		return validComparableRequired[uint](offset, field)

	case reflect.Uint8:
		return validComparableRequired[uint8](offset, field)

	case reflect.Uint16:
		return validComparableRequired[uint16](offset, field)

	case reflect.Uint32:
		return validComparableRequired[uint32](offset, field)

	case reflect.Uint64:
		return validComparableRequired[uint64](offset, field)

	case reflect.Float32:
		return validComparableRequired[float32](offset, field)

	case reflect.Float64:
		return validComparableRequired[float64](offset, field)

	// case reflect.Array:
	// case reflect.Slice:

	case reflect.String:
		return validComparableRequired[string](offset, field)

	default:
		return nil, notImplemented("required", kind)
	}
}

func validComparableRequired[T comparable](offset uintptr, field string) (validator, error) {
	var zero T

	return func(ptr unsafe.Pointer, errs *FieldErrors) {
		if *(*T)(unsafe.Add(ptr, offset)) == zero {
			errs.Append(FieldError{
				err:    ErrRequired,
				field:  field,
				expect: "any value",
			})
		}
	}, nil
}
