package valid

import (
	"reflect"
	"unsafe"
)

type validator func(ptr unsafe.Pointer, errs *FieldErrors)
type validatorCreator func(offset uintptr, typ reflect.Type, field string, s string) (validator, error)

// func validNumMax[T constraints.Number](s string) (validator, error) {
// 	var max T

// 	if err := scanner.ScanString(&max, s); err != nil {
// 		return nil, err
// 	}

// 	var e = Error("value can't exceed %d", max)

// 	return func(ptr unsafe.Pointer) error {
// 		if *(*T)(ptr) > max {
// 			return e
// 		}

// 		return nil
// 	}, nil
// }

// func validStrMin(s string) (validator, error) {
// 	var min int

// 	if err := scanner.ScanString(&min, s); err != nil {
// 		return nil, err
// 	}

// 	var e = Error("length must be at least %d", min)

// 	return func(ptr unsafe.Pointer) error {
// 		if len(*(*string)(ptr)) < min {
// 			return e
// 		}

// 		return nil
// 	}, nil
// }

// func validStrMax(s string) (validator, error) {
// 	var max int

// 	if err := scanner.ScanString(&max, s); err != nil {
// 		return nil, err
// 	}

// 	var e = Error("length can't exceed %d", max)

// 	return func(ptr unsafe.Pointer) error {
// 		if len(*(*string)(ptr)) > max {
// 			return e
// 		}

// 		return nil
// 	}, nil
// }
