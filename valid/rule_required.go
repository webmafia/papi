package valid

import (
	"reflect"
	"unsafe"
)

func appendRequiredValidators(valids *validators, offset uintptr, fld *reflect.StructField) (err error) {
	switch kind := fld.Type.Kind(); kind {

	case reflect.Int:
		return valids.append(validComparableRequired[int](offset, fld.Name))

	case reflect.Int8:
		return valids.append(validComparableRequired[int8](offset, fld.Name))

	case reflect.Int16:
		return valids.append(validComparableRequired[int16](offset, fld.Name))

	case reflect.Int32:
		return valids.append(validComparableRequired[int32](offset, fld.Name))

	case reflect.Int64:
		return valids.append(validComparableRequired[int64](offset, fld.Name))

	case reflect.Uint:
		return valids.append(validComparableRequired[uint](offset, fld.Name))

	case reflect.Uint8:
		return valids.append(validComparableRequired[uint8](offset, fld.Name))

	case reflect.Uint16:
		return valids.append(validComparableRequired[uint16](offset, fld.Name))

	case reflect.Uint32:
		return valids.append(validComparableRequired[uint32](offset, fld.Name))

	case reflect.Uint64:
		return valids.append(validComparableRequired[uint64](offset, fld.Name))

	case reflect.Float32:
		return valids.append(validComparableRequired[float32](offset, fld.Name))

	case reflect.Float64:
		return valids.append(validComparableRequired[float64](offset, fld.Name))

	// case reflect.Array:
	// case reflect.Slice:

	case reflect.String:
		return valids.append(validComparableRequired[string](offset, fld.Name))
	}

	return
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
