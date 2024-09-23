package strings

import (
	"reflect"
	"unsafe"

	"github.com/webmafia/fast"
	"github.com/webmafia/fastapi/internal"
)

func (f *Factory) createArrayScanner(typ reflect.Type) (scan Scanner, err error) {
	const sep byte = ','

	elem := typ.Elem()
	arrSize := typ.Len()
	itemSize := elem.Size()
	elemScan, err := f.Scanner(elem)

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