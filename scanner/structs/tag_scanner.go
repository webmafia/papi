package structs

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/webmafia/fastapi/scanner/strings"
)

func CreateTagScanner(f *strings.Factory, typ reflect.Type) (scan strings.Scanner, err error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("invalid struct")
	}

	numFields := typ.NumField()

	type field struct {
		scan   strings.Scanner
		offset uintptr
	}

	var fldScan strings.Scanner

	tagScanners := make(map[string]field, numFields)

	for i := range numFields {
		fld := typ.Field(i)

		for k, v := range iterateStructTags(fld.Tag) {
			if k != "tag" {
				continue
			}

			if fldScan, err = f.Scanner(fld.Type); err != nil {
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
