package scanner

import (
	"reflect"
	"unsafe"
)

func (c Creator) createPointerScanner(typ reflect.Type) (scan Scanner, err error) {
	elem := typ.Elem()
	elemScan, err := c.CreateScanner(elem)

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
