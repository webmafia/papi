package valid

import (
	"reflect"
	"slices"
	"unsafe"

	"github.com/webbmaffian/papi/registry/scanner"
)

func appendEnumValidators(valids *validators, offset uintptr, fld *reflect.StructField, s string) (err error) {
	switch kind := fld.Type.Kind(); kind {

	case reflect.Int:
		return valids.append(validComparableEnum[int](offset, fld.Name, s))

	case reflect.Int8:
		return valids.append(validComparableEnum[int8](offset, fld.Name, s))

	case reflect.Int16:
		return valids.append(validComparableEnum[int16](offset, fld.Name, s))

	case reflect.Int32:
		return valids.append(validComparableEnum[int32](offset, fld.Name, s))

	case reflect.Int64:
		return valids.append(validComparableEnum[int64](offset, fld.Name, s))

	case reflect.Uint:
		return valids.append(validComparableEnum[uint](offset, fld.Name, s))

	case reflect.Uint8:
		return valids.append(validComparableEnum[uint8](offset, fld.Name, s))

	case reflect.Uint16:
		return valids.append(validComparableEnum[uint16](offset, fld.Name, s))

	case reflect.Uint32:
		return valids.append(validComparableEnum[uint32](offset, fld.Name, s))

	case reflect.Uint64:
		return valids.append(validComparableEnum[uint64](offset, fld.Name, s))

	case reflect.Float32:
		return valids.append(validComparableEnum[float32](offset, fld.Name, s))

	case reflect.Float64:
		return valids.append(validComparableEnum[float64](offset, fld.Name, s))

	// case reflect.Array:
	// case reflect.Slice:

	case reflect.String:
		return valids.append(validComparableEnum[string](offset, fld.Name, s))
	}

	return
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
