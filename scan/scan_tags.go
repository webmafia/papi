package scan

import (
	"fmt"
	"reflect"
	"unsafe"
)

func ScanTags[T any](v *T, tags reflect.StructTag) (err error) {
	typ := reflect.TypeOf(*v)

	if kind := typ.Kind(); kind != reflect.Struct {
		return fmt.Errorf("expected struct, but got %s", kind)
	}

	return scanStructTags(unsafe.Pointer(v), typ, tags)
}

func scanStructTags(ptr unsafe.Pointer, typ reflect.Type, tags reflect.StructTag) (err error) {
	numFlds := typ.NumField()

	for i := 0; i < numFlds; i++ {
		fld := typ.Field(i)

		if fld.Type.Kind() == reflect.Struct {
			if err = scanStructTags(unsafe.Add(ptr, fld.Offset), fld.Type, tags); err != nil {
				return
			}

			continue
		}

		tag, ok := fld.Tag.Lookup("tag")

		if !ok {
			continue
		}

		val, ok := tags.Lookup(tag)

		if !ok {
			continue
		}

		scan, err := CreateScanner(fld.Type)

		if err != nil {
			return err
		}

		if err := scan(unsafe.Add(ptr, fld.Offset), val); err != nil {
			return err
		}
	}

	return
}
