package internal

import (
	"errors"
	"reflect"
	"unsafe"
)

type StructTagScanner func(dst unsafe.Pointer, src reflect.StructTag) error

func CreateStructTagScanner(typ reflect.Type) (scan StructTagScanner, err error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("invalid struct")
	}

	numFields := typ.NumField()

	type field struct {
		scan   Scanner
		offset uintptr
	}

	var fldScan Scanner

	tagScanners := make(map[string]field, numFields)

	for i := range numFields {
		fld := typ.Field(i)

		for k, v := range iterateStructTags(fld.Tag) {
			if k != "tag" {
				continue
			}

			if fldScan, err = CreateScanner(fld.Type); err != nil {
				return
			}

			tagScanners[v] = field{
				scan:   fldScan,
				offset: fld.Offset,
			}

			break
		}
	}

	return func(dst unsafe.Pointer, src reflect.StructTag) (err error) {
		for k, v := range iterateStructTags(src) {
			if fld, ok := tagScanners[k]; ok {
				if err = fld.scan(unsafe.Add(dst, fld.offset), v); err != nil {
					return
				}
			}
		}

		return
	}, nil
}
