package valid

import (
	"fmt"
	"reflect"
	"unsafe"
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

func stringLen(ptr unsafe.Pointer) int {
	return (*stringHeader)(ptr).len
}

func notImplemented(validation string, kind reflect.Kind) error {
	return fmt.Errorf("'%s' validation of %s is not implemented", validation, kind)
}
