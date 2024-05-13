package datatype

import (
	"fmt"
	"reflect"
	"unsafe"
)

func Scan[T any](d *DataTypes, v *T, str string) (err error) {
	typ := reflect.TypeOf(v)
	scan, ok := d.scanners[typ]

	if !ok {
		return fmt.Errorf("missing scanner for type: %s", typ)
	}

	return scan(unsafe.Pointer(v), str)
}
