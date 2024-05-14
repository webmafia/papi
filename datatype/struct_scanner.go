package datatype

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/webmafia/fast"
)

type structField struct {
	tag    string
	val    string
	scan   func(unsafe.Pointer, string) error
	offset uintptr
}

func CreateStructScanner[T any](d *DataTypes, tags ...string) (fn func(*T, func(tag, val string) string) error, err error) {
	typ := reflect.TypeOf((*T)(nil)).Elem()

	if kind := typ.Kind(); kind != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", kind)
	}

	numFlds := typ.NumField()
	flds := make([]structField, 0, numFlds)

	for i := 0; i < numFlds; i++ {
		fld := typ.Field(i)

		if len(tags) > 0 {
			for _, tag := range tags {
				val, ok := fld.Tag.Lookup(tag)

				if !ok {
					continue
				}

				scan, ok := d.scanners[reflect.PointerTo(fld.Type)]

				if !ok {
					return nil, fmt.Errorf("no scanner found for %s", fld.Type)
				}

				flds = append(flds, structField{
					tag:    tag,
					val:    val,
					scan:   scan,
					offset: fld.Offset,
				})
			}
		} else {
			scan, ok := d.scanners[reflect.PointerTo(fld.Type)]

			if !ok {
				return nil, fmt.Errorf("no scanner found for %s", fld.Type)
			}

			flds = append(flds, structField{
				val:    fld.Name,
				scan:   scan,
				offset: fld.Offset,
			})
		}
	}

	// Save memory by removing any capacity leftovers
	if cap(flds) > len(flds) {
		newFlds := make([]structField, len(flds))
		copy(newFlds, flds)
		flds = newFlds
	}

	fn = func(v *T, cb func(tag, val string) string) (err error) {
		ptr := fast.Noescape(unsafe.Pointer(v))

		for i := range flds {
			str := cb(flds[i].tag, flds[i].val)

			if str == "" {
				continue
			}

			if err = flds[i].scan(unsafe.Add(ptr, flds[i].offset), str); err != nil {
				return
			}
		}

		return
	}

	return
}
