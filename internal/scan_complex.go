package internal

import (
	"strconv"
	"unsafe"
)

func scanComplex64(p unsafe.Pointer, s string) (err error) {
	v, err := strconv.ParseComplex(s, 64)

	if err == nil {
		*(*complex64)(p) = complex64(v)
	}

	return
}

func scanComplex128(p unsafe.Pointer, s string) (err error) {
	v, err := strconv.ParseComplex(s, 128)

	if err == nil {
		*(*complex128)(p) = v
	}

	return
}
