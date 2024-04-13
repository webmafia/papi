package scan

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

func CreateStructScanner(typ reflect.Type) (scan *StructScannerBuilder, err error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("invalid struct")
	}

	return &StructScannerBuilder{
		typ:       typ,
		numFields: typ.NumField(),
	}, nil
}

type StructScannerBuilder struct {
	typ       reflect.Type
	numFields int
	fields    []fieldScanner
}

type fieldScanner struct {
	offset uintptr
	scan   Scanner
}

func (b *StructScannerBuilder) AddByTag(tag string, val string) (err error) {
	for i := 0; i < b.numFields; i++ {
		fld := b.typ.Field(i)
		tagVal := fld.Tag.Get(tag)

		if tagVal != val {
			continue
		}

		return b.AddField(fld)
	}

	return fmt.Errorf("tag '%s' of value '%s' not found", tag, val)
}

func (b *StructScannerBuilder) AddField(fld reflect.StructField) (err error) {
	scan, err := CreateScanner(fld.Type)

	if err != nil {
		return err
	}

	b.fields = append(b.fields, fieldScanner{
		offset: fld.Offset,
		scan:   scan,
	})

	return
}

func (b *StructScannerBuilder) Compile() StructScanner {
	flds := b.fields

	return func(p unsafe.Pointer, s ...string) (err error) {
		l := min(len(flds), len(s))

		for i := 0; i < l; i++ {
			if err = flds[i].scan(unsafe.Add(p, flds[i].offset), s[i]); err != nil {
				return
			}
		}

		return
	}
}

type StructScanner func(unsafe.Pointer, ...string) error
