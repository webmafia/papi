package scanner

import (
	"reflect"
	"unsafe"

	"github.com/webbmaffian/papi/internal/iterate"
	"github.com/webmafia/fast"
)

type sliceHeader struct {
	data unsafe.Pointer
	len  int
	cap  int
}

func (c Creator) createSliceScanner(typ reflect.Type) (scan Scanner, err error) {
	const sep byte = ','

	elem := typ.Elem()
	itemSize := elem.Size()
	elemScan, err := c.CreateScanner(elem)

	if err != nil {
		return
	}

	scan = func(p unsafe.Pointer, s string) (err error) {
		head := (*sliceHeader)(p)
		var calcSize int

		for range iterate.IterateChunks(s, sep) {
			calcSize++
		}

		if calcSize > head.cap {

			// Allocate new slice with the calculated size
			newSlice := reflect.MakeSlice(typ, calcSize, calcSize)
			head.data = newSlice.UnsafePointer()
			head.cap = calcSize
			head.len = 0
		}

		for i, sub := range iterate.IterateChunks(s, sep) {
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
