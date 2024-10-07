package value

import (
	"reflect"
	"unsafe"
)

func createPointerScanner(typ reflect.Type, createElemScanner CreateValueScanner) (scan ValueScanner, err error) {
	elem := typ.Elem()
	elemScan, err := createElemScanner(elem, createElemScanner)

	if err != nil {
		return
	}

	scan = func(p unsafe.Pointer, s string) error {
		subPtr := (*unsafe.Pointer)(p)

		if *subPtr == nil {
			*subPtr = reflect.New(elem).UnsafePointer()
		}

		return elemScan(*subPtr, s)
	}

	return
}
