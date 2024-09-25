package value

import (
	"reflect"
	"unsafe"

	"github.com/webmafia/fast"
	"github.com/webmafia/fastapi/internal"
)

func createArrayScanner(typ reflect.Type, createElemScanner CreateValueScanner) (scan ValueScanner, err error) {
	const sep byte = ','

	elem := typ.Elem()
	arrSize := typ.Len()
	itemSize := elem.Size()
	elemScan, err := createElemScanner(elem, createElemScanner)

	if err != nil {
		return
	}

	scan = func(p unsafe.Pointer, s string) (err error) {
		for i, sub := range internal.IterateChunks(s, sep) {
			if i >= arrSize {
				break
			}

			elemPtr := unsafe.Add(p, uintptr(i)*itemSize)

			if err = elemScan(fast.Noescape(elemPtr), sub); err != nil {
				return err
			}
		}

		return
	}

	return
}
