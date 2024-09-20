package fastapi

import (
	"strconv"
	"unsafe"
)

func scanUint(p unsafe.Pointer, s string) (err error) {
	v, err := strconv.ParseUint(s, 10, 0)

	if err == nil {
		*(*uint)(p) = uint(v)
	}

	return
}

func scanUint8(p unsafe.Pointer, s string) (err error) {
	v, err := strconv.ParseUint(s, 10, 8)

	if err == nil {
		*(*uint8)(p) = uint8(v)
	}

	return
}

func scanUint16(p unsafe.Pointer, s string) (err error) {
	v, err := strconv.ParseUint(s, 10, 16)

	if err == nil {
		*(*uint16)(p) = uint16(v)
	}

	return
}

func scanUint32(p unsafe.Pointer, s string) (err error) {
	v, err := strconv.ParseUint(s, 10, 32)

	if err == nil {
		*(*uint32)(p) = uint32(v)
	}

	return
}

func scanUint64(p unsafe.Pointer, s string) (err error) {
	v, err := strconv.ParseUint(s, 10, 64)

	if err == nil {
		*(*uint64)(p) = v
	}

	return
}
