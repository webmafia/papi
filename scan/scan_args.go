package scan

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
)

func CreateArgsScanner(typ reflect.Type) (scan func(unsafe.Pointer, *fasthttp.Args) error, err error) {
	type field struct {
		key    string
		offset uintptr
		scan   Scanner
	}

	var fields []field
	l := typ.NumField()

	for i := 0; i < l; i++ {
		fld := typ.Field(i)

		if key, ok := fld.Tag.Lookup("query"); ok {
			scan, err := CreateScanner(fld.Type)

			if err != nil {
				return nil, err
			}

			fields = append(fields, field{
				key:    key,
				offset: fld.Offset,
				scan:   scan,
			})
		}
	}

	return func(p unsafe.Pointer, args *fasthttp.Args) error {
		for _, fld := range fields {
			val := args.Peek(fld.key)

			if len(val) > 0 {
				if err := fld.scan(unsafe.Add(p, fld.offset), fast.BytesToString(val)); err != nil {
					return err
				}
			}
		}

		return nil
	}, nil
}
