package scanner

import (
	"reflect"
	"unsafe"

	"github.com/webmafia/fast"
	"github.com/webmafia/papi/internal/iterate"
)

func (c Creator) createArrayScanner(typ reflect.Type) (scan Scanner, err error) {
	const sep byte = ','

	elem := typ.Elem()
	arrSize := typ.Len()
	itemSize := elem.Size()
	elemScan, err := c.CreateScanner(elem)

	if err != nil {
		return
	}

	scan = func(p unsafe.Pointer, s string) (err error) {
		for i, sub := range iterate.IterateChunks(s, sep) {
			if i >= arrSize {
				break
			}

			elemPtr := unsafe.Add(p, uintptr(i)*itemSize)

			if err = elemScan(fast.NoescapeUnsafe(elemPtr), sub); err != nil {
				return err
			}
		}

		return
	}

	return
}
