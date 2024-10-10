package valid

import "unsafe"

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
