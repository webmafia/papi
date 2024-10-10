package valid

import (
	"reflect"
	"unsafe"

	"github.com/webbmaffian/papi/internal/constraints"
	"github.com/webbmaffian/papi/registry/scanner"
)

func appendMinValidators(valids *validators, offset uintptr, fld *reflect.StructField, s string) (err error) {
	switch kind := fld.Type.Kind(); kind {

	case reflect.Int:
		return valids.append(validNumMin[int](offset, fld.Name, s))

	case reflect.Int8:
		return valids.append(validNumMin[int8](offset, fld.Name, s))

	case reflect.Int16:
		return valids.append(validNumMin[int16](offset, fld.Name, s))

	case reflect.Int32:
		return valids.append(validNumMin[int32](offset, fld.Name, s))

	case reflect.Int64:
		return valids.append(validNumMin[int64](offset, fld.Name, s))

	case reflect.Uint:
		return valids.append(validNumMin[uint](offset, fld.Name, s))

	case reflect.Uint8:
		return valids.append(validNumMin[uint8](offset, fld.Name, s))

	case reflect.Uint16:
		return valids.append(validNumMin[uint16](offset, fld.Name, s))

	case reflect.Uint32:
		return valids.append(validNumMin[uint32](offset, fld.Name, s))

	case reflect.Uint64:
		return valids.append(validNumMin[uint64](offset, fld.Name, s))

	case reflect.Float32:
		return valids.append(validNumMin[float32](offset, fld.Name, s))

	case reflect.Float64:
		return valids.append(validNumMin[float64](offset, fld.Name, s))

	case reflect.Slice:
		return valids.append(validSliceMin(offset, fld.Name, s))

	case reflect.String:
		return valids.append(validStringMin(offset, fld.Name, s))
	}

	return
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
