package strings

import (
	"reflect"
	"unsafe"
)

type Scanner func(unsafe.Pointer, string) error

func RegisterScanner[T any](f *Factory, fn func(dst *T, src string) error) {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	f.RegisterScanner(typ, *(*Scanner)(unsafe.Pointer(&fn)))
}

func GetScanner[T any](f *Factory, dst *T) (fn func(*T, string) error, err error) {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	scan, err := f.Scanner(typ)

	if err != nil {
		return
	}

	fn = *(*func(*T, string) error)(unsafe.Pointer(&scan))
	return
}
