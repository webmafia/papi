package valid

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/webbmaffian/papi/errors"
)

type sliceHeader struct {
	data unsafe.Pointer
	cap  int
	len  int
}

type stringHeader struct {
	data unsafe.Pointer
	len  int
}

func sliceLen(ptr unsafe.Pointer) int {
	return (*sliceHeader)(ptr).len
}

func sliceDataAndLen(ptr unsafe.Pointer) (unsafe.Pointer, int) {
	return (*sliceHeader)(ptr).data, (*sliceHeader)(ptr).len
}

func stringLen(ptr unsafe.Pointer) int {
	return (*stringHeader)(ptr).len
}

func notImplemented(validation string, kind reflect.Kind) error {
	return fmt.Errorf("'%s' validation of %s is not implemented", validation, kind)
}

func validArray(offset uintptr, typ reflect.Type, field string, s string, create validatorCreator) (validator, error) {
	elem := typ.Elem()
	valid, err := create(offset, elem, field, s)

	if err != nil {
		return nil, err
	}

	l := typ.Len()
	size := elem.Size()

	return func(ptr unsafe.Pointer, errs *errors.Errors) {
		for i := range l {
			valid(unsafe.Add(ptr, uintptr(i)*size), errs)
		}
	}, nil
}

func validSlice(offset uintptr, typ reflect.Type, field string, s string, create validatorCreator) (validator, error) {
	elem := typ.Elem()
	valid, err := create(0, elem, field, s)

	if err != nil {
		return nil, err
	}

	size := elem.Size()

	return func(ptr unsafe.Pointer, errs *errors.Errors) {
		data, l := sliceDataAndLen(unsafe.Add(ptr, offset))

		for i := range l {
			valid(unsafe.Add(data, uintptr(i)*size), errs)
		}
	}, nil
}

func validPointer(offset uintptr, typ reflect.Type, field string, s string, create validatorCreator) (validator, error) {
	elem := typ.Elem()
	valid, err := create(0, elem, field, s)

	if err != nil {
		return nil, err
	}

	return func(ptr unsafe.Pointer, errs *errors.Errors) {
		valid(*(*unsafe.Pointer)(unsafe.Add(ptr, offset)), errs)
	}, nil
}
