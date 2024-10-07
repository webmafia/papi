package structs

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/webmafia/fastapi/registry"
	"github.com/webmafia/fastapi/registry/value"
)

func CreateTagScanner(r *registry.Registry, typ reflect.Type) (scan value.ValueScanner, err error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("invalid struct")
	}

	numFields := typ.NumField()

	type field struct {
		scan   value.ValueScanner
		offset uintptr
	}

	var fldScan value.ValueScanner

	tagScanners := make(map[string]field, numFields)

	for i := range numFields {
		fld := typ.Field(i)

		for k, v := range iterateStructTags(fld.Tag) {
			if k != "tag" {
				continue
			}

			if fldScan, err = r.CreateValueScanner(fld.Type, fld.Tag); err != nil {
				return
			}

			tagScanners[v] = field{
				scan:   fldScan,
				offset: fld.Offset,
			}

			break
		}
	}

	return func(dst unsafe.Pointer, src string) (err error) {
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
