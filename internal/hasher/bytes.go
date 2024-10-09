package hasher

import (
	"unsafe"

	"github.com/webmafia/fast"
)

//go:inline
func toBytes[T any](v *T) []byte {
	return fast.PointerToBytes(v, int(unsafe.Sizeof(*v)))
}
