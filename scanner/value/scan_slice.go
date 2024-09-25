package value

import (
	"reflect"
	"unsafe"

	"github.com/webmafia/fast"
	"github.com/webmafia/fastapi/internal"
)

type sliceHeader struct {
	data unsafe.Pointer
	len  int
	cap  int
}

func createSliceScanner(typ reflect.Type, createElemScanner CreateValueScanner) (scan ValueScanner, err error) {
	const sep byte = ','

	elem := typ.Elem()
	itemSize := elem.Size()
	elemScan, err := createElemScanner(elem, createElemScanner)

	if err != nil {
		return
	}

	scan = func(p unsafe.Pointer, s string) (err error) {
		head := (*sliceHeader)(p)
		var calcSize int

		for range internal.IterateChunks(s, sep) {
			calcSize++
		}

		if calcSize > head.cap {

			// Allocate new slice with the calculated size
			newSlice := reflect.MakeSlice(typ, calcSize, calcSize)
			head.data = newSlice.UnsafePointer()
			head.cap = calcSize
			head.len = 0
		}

		for i, sub := range internal.IterateChunks(s, sep) {
			elemPtr := unsafe.Add(head.data, uintptr(i)*itemSize)

			if err = elemScan(fast.Noescape(elemPtr), sub); err != nil {
				return err
			}
		}

		head.len = calcSize
		return
	}

	return
}
