package strings

import (
	"reflect"
	"unsafe"

	"github.com/webmafia/fast"
)

func ScanString[T any](f *Factory, dst *T, src string) (err error) {
	scan, err := f.Scanner(reflect.TypeOf(dst).Elem())

	if err != nil {
		return
	}

	return scan(fast.Noescape(unsafe.Pointer(dst)), src)
}

func ScanBytes[T any](f *Factory, dst *T, src []byte) (err error) {
	return ScanString(f, dst, fast.BytesToString(src))
}
