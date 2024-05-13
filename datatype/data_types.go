package datatype

import (
	"reflect"
	"unsafe"
)

type DataTypes struct {
	scanners map[reflect.Type]func(unsafe.Pointer, string) error
	encoders map[reflect.Type]func(unsafe.Pointer, string) error
}

func NewDataTypes() *DataTypes {
	return &DataTypes{
		scanners: make(map[reflect.Type]func(unsafe.Pointer, string) error),
		encoders: make(map[reflect.Type]func(unsafe.Pointer, string) error),
	}
}
