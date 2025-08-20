package registry

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"unsafe"

	"github.com/webmafia/papi/internal/iterate"
)

func ScanTags[T any](reg *Registry, dst *T, tags reflect.StructTag) (err error) {
	typ := reflect.TypeFor[T]()
	scan, ok := reg.scanCache[typ]

	if !ok {
		if scan, err = reg.createTagScanner(typ); err != nil {
			return
		}

		reg.scanCache[typ] = scan
	}

	return scan(unsafe.Pointer(dst), string(tags))
}

func scanSchemaTags[T any](reg *Registry, dst *T, tags reflect.StructTag) (*T, error) {
	if err := ScanTags(reg, dst, tags); err != nil {
		return nil, err
	}

	return dst, nil
}

func (r *Registry) createTagScanner(typ reflect.Type) (scan Parser, err error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("invalid struct")
	}

	numFields := typ.NumField()

	type field struct {
		scan   Parser
		offset uintptr
	}

	var fldScan Parser

	tagScanners := make(map[string]field, numFields)

	var (
		flags       []string
		flagOffsets []uintptr
	)

	for i := range numFields {
		fld := typ.Field(i)

		for k, v := range iterate.IterateStructTags(fld.Tag) {
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

			if fldScan, err = r.Parser(fld.Type, fld.Tag); err != nil {
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
		for k, v := range iterate.IterateStructTags(src) {
			if k == "flags" {
				for flag := range iterate.IterateFlags(v) {
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
