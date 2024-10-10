package structs

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"unsafe"

	"github.com/webbmaffian/papi/registry/scanner"
)

func CreateTagScanner(typ reflect.Type, createValueScanner func(typ reflect.Type, tags reflect.StructTag) (scan scanner.Scanner, err error)) (scan scanner.Scanner, err error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("invalid struct")
	}

	numFields := typ.NumField()

	type field struct {
		scan   scanner.Scanner
		offset uintptr
	}

	var fldScan scanner.Scanner

	tagScanners := make(map[string]field, numFields)

	var (
		flags       []string
		flagOffsets []uintptr
	)

	for i := range numFields {
		fld := typ.Field(i)

		for k, v := range IterateStructTags(fld.Tag) {
			if k != "tag" {
				continue
			}

			if fld.Type.Kind() == reflect.Bool {
				if k, v, ok := strings.Cut(v, ":"); ok && k == "flags" {
					flags = append(flags, v)
					flagOffsets = append(flagOffsets, fld.Offset)
					continue
				}
			}

			if fldScan, err = createValueScanner(fld.Type, fld.Tag); err != nil {
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
		for k, v := range IterateStructTags(src) {
			if k == "flags" {
				for flag := range iterateFlags(v) {
					idx := slices.Index(flags, flag)

					if idx < 0 {
						continue
					}

					offset := flagOffsets[idx]
					*(*bool)(unsafe.Add(dst, offset)) = true
				}
			} else if fld, ok := tagScanners[k]; ok {
				if err = fld.scan(unsafe.Add(dst, fld.offset), v); err != nil {
					return
				}
			}
		}

		return
	}, nil
}
