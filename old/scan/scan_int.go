package scan

import (
	"strconv"
	"unsafe"
)

func scanInt(p unsafe.Pointer, s string) (err error) {
	v, err := strconv.ParseInt(s, 10, 0)

	if err == nil {
		*(*int)(p) = int(v)
	}

	return
}

func scanInt8(p unsafe.Pointer, s string) (err error) {
	v, err := strconv.ParseInt(s, 10, 8)

	if err == nil {
		*(*int8)(p) = int8(v)
	}

	return
}

func scanInt16(p unsafe.Pointer, s string) (err error) {
	v, err := strconv.ParseInt(s, 10, 16)

	if err == nil {
		*(*int16)(p) = int16(v)
	}

	return
}

func scanInt32(p unsafe.Pointer, s string) (err error) {
	v, err := strconv.ParseInt(s, 10, 32)

	if err == nil {
		*(*int32)(p) = int32(v)
	}

	return
}

func scanInt64(p unsafe.Pointer, s string) (err error) {
	v, err := strconv.ParseInt(s, 10, 64)

	if err == nil {
		*(*int64)(p) = v
	}

	return
}
