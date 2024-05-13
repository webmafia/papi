package datatype

import (
	"reflect"
	"unsafe"
)

func RegisterScanner[T any](d *DataTypes, fn func(ptr *T, str string) (err error)) {
	d.registerScanner(
		reflect.TypeOf((*T)(nil)),
		*(*func(unsafe.Pointer, string) error)(unsafe.Pointer(&fn)),
	)
}

func (d *DataTypes) registerScanner(typ reflect.Type, fn func(unsafe.Pointer, string) error) {
	d.scanners[typ] = fn
}
