package fastapi

import (
	"fmt"
	"unsafe"
)

func scanBool(p unsafe.Pointer, s string) (err error) {
	var v bool

	switch s {
	case "1", "t", "T", "true", "TRUE", "True", "yes", "YES", "Yes":
		v = true
	case "0", "f", "F", "false", "FALSE", "False", "no", "NO", "No":
		v = false
	default:
		return fmt.Errorf("invalid boolean: '%s'", s)
	}

	*(*bool)(p) = v

	return
}
