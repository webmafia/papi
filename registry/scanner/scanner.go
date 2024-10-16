package scanner

import (
	"unsafe"
)

type Scanner func(unsafe.Pointer, string) error
