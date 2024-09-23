package strings

import (
	"reflect"
	"unsafe"
)

func (f *Factory) createPointerScanner(typ reflect.Type) (scan Scanner, err error) {
	elem := typ.Elem()
	elemScan, err := f.Scanner(elem)

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
