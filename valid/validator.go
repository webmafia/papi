package valid

import (
	"unsafe"

	"github.com/webbmaffian/papi/internal/constraints"
	"github.com/webbmaffian/papi/registry/scanner"
)

type validator func(ptr unsafe.Pointer) error

func validNumMin[T constraints.Number](offset uintptr, name string, s string) (validator, error) {
	var min T

	if err := scanner.ScanString(&min, s); err != nil {
		return nil, err
	}

	var e = Error("%s: value can't be less than %d", name, min)

	return func(ptr unsafe.Pointer) error {
		if *(*T)(unsafe.Add(ptr, offset)) < min {
			return e
		}

		return nil
	}, nil
}

func validNumMax[T constraints.Number](s string) (validator, error) {
	var max T

	if err := scanner.ScanString(&max, s); err != nil {
		return nil, err
	}

	var e = Error("value can't exceed %d", max)

	return func(ptr unsafe.Pointer) error {
		if *(*T)(ptr) > max {
			return e
		}

		return nil
	}, nil
}

func validStrMin(s string) (validator, error) {
	var min int

	if err := scanner.ScanString(&min, s); err != nil {
		return nil, err
	}

	var e = Error("length must be at least %d", min)

	return func(ptr unsafe.Pointer) error {
		if len(*(*string)(ptr)) < min {
			return e
		}

		return nil
	}, nil
}

func validStrMax(s string) (validator, error) {
	var max int

	if err := scanner.ScanString(&max, s); err != nil {
		return nil, err
	}

	var e = Error("length can't exceed %d", max)

	return func(ptr unsafe.Pointer) error {
		if len(*(*string)(ptr)) > max {
			return e
		}

		return nil
	}, nil
}
