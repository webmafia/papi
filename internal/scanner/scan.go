package scanner

import (
	"reflect"
	"unsafe"

	"github.com/webmafia/fast"
)

var c = NewCreator()

func getScanner[T any](dst *T) (scan Scanner, err error) {
	return c.CreateScanner(reflect.TypeOf(dst).Elem())
}

//go:inline
func CreateScanner(typ reflect.Type) (scan Scanner, err error) {
	return c.CreateScanner(typ)
}

func GetScanner[T any](dst *T) (scan func(*T, string) error, err error) {
	sc, err := getScanner(dst)

	if err != nil {
		return
	}

	return *(*func(*T, string) error)(unsafe.Pointer(&sc)), nil
}

func ScanString[T any](dst *T, src string) (err error) {
	scan, err := getScanner(dst)

	if err != nil {
		return
	}

	return scan(fast.Noescape(unsafe.Pointer(dst)), src)
}

func ScanBytes[T any](dst *T, src []byte) (err error) {
	return ScanString(dst, fast.BytesToString(src))
}
