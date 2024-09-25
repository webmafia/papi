package value

import (
	"strconv"
	"unsafe"
)

func scanFloat32(p unsafe.Pointer, s string) (err error) {
	v, err := strconv.ParseFloat(s, 32)

	if err == nil {
		*(*float32)(p) = float32(v)
	}

	return
}

func scanFloat64(p unsafe.Pointer, s string) (err error) {
	v, err := strconv.ParseFloat(s, 64)

	if err == nil {
		*(*float64)(p) = v
	}

	return
}
