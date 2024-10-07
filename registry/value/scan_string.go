package value

import (
	"strings"
	"unsafe"
)

func scanString(p unsafe.Pointer, s string) (err error) {
	*(*string)(p) = s
	return
}

func scanStringCopy(p unsafe.Pointer, s string) (err error) {
	*(*string)(p) = strings.Clone(s)
	return
}
