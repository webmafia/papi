package valid

import (
	"reflect"
	"unsafe"

	"github.com/webbmaffian/papi/internal/constraints"
	"github.com/webbmaffian/papi/registry/scanner"
)

func appendMaxValidators(valids *validators, offset uintptr, fld *reflect.StructField, s string) (err error) {
	switch kind := fld.Type.Kind(); kind {

	case reflect.Int:
		return valids.append(validNumMax[int](offset, fld.Name, s))

	case reflect.Int8:
		return valids.append(validNumMax[int8](offset, fld.Name, s))

	case reflect.Int16:
		return valids.append(validNumMax[int16](offset, fld.Name, s))

	case reflect.Int32:
		return valids.append(validNumMax[int32](offset, fld.Name, s))

	case reflect.Int64:
		return valids.append(validNumMax[int64](offset, fld.Name, s))

	case reflect.Uint:
		return valids.append(validNumMax[uint](offset, fld.Name, s))

	case reflect.Uint8:
		return valids.append(validNumMax[uint8](offset, fld.Name, s))

	case reflect.Uint16:
		return valids.append(validNumMax[uint16](offset, fld.Name, s))

	case reflect.Uint32:
		return valids.append(validNumMax[uint32](offset, fld.Name, s))

	case reflect.Uint64:
		return valids.append(validNumMax[uint64](offset, fld.Name, s))

	case reflect.Float32:
		return valids.append(validNumMax[float32](offset, fld.Name, s))

	case reflect.Float64:
		return valids.append(validNumMax[float64](offset, fld.Name, s))

	case reflect.Slice:
		return valids.append(validSliceMax(offset, fld.Name, s))

	case reflect.String:
		return valids.append(validStringMax(offset, fld.Name, s))
	}

	return
}

func validNumMax[T constraints.Number](offset uintptr, field string, s string) (validator, error) {
	var max T

	if err := scanner.ScanString(&max, s); err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *FieldErrors) {
		if val := *(*T)(unsafe.Add(ptr, offset)); val != 0 && val > max {
			errs.Append(FieldError{
				err:    ErrAboveMax,
				field:  field,
				expect: s,
			})
		}
	}, nil
}

func validSliceMax(offset uintptr, field string, s string) (validator, error) {
	var max int

	if err := scanner.ScanString(&max, s); err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *FieldErrors) {
		if l := sliceLen(unsafe.Add(ptr, offset)); l != 0 && l > max {
			errs.Append(FieldError{
				err:    ErrAboveMax,
				field:  field,
				expect: s,
			})
		}
	}, nil
}

func validStringMax(offset uintptr, field string, s string) (validator, error) {
	var max int

	if err := scanner.ScanString(&max, s); err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *FieldErrors) {
		if l := stringLen(unsafe.Add(ptr, offset)); l != 0 && l > max {
			errs.Append(FieldError{
				err:    ErrAboveMax,
				field:  field,
				expect: s,
			})
		}
	}, nil
}
