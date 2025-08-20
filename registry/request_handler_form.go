package registry

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
)

type formDecoder struct {
	offset uintptr
	scan   Decoder
	name   string
}

func (r *Registry) createFormHandler(typ reflect.Type) (scan Handler, err error) {
	if typ.Kind() != reflect.Struct {
		err = errors.New("invalid struct")
		return
	}

	numFields := typ.NumField()
	dec := make([]formDecoder, 0, numFields)

	for i := 0; i < numFields; i++ {
		fld := typ.Field(i)
		name := fld.Tag.Get("form")

		if name == "" {
			name = fld.Name
		}

		if name == "-" {
			continue
		}

		sc, err := r.Decoder(fld.Type, fld.Tag)

		if err != nil {
			return nil, err
		}

		dec = append(dec, formDecoder{
			offset: fld.Offset,
			scan:   sc,
			name:   name,
		})
	}

	scan = func(c *fasthttp.RequestCtx, p unsafe.Pointer) (err error) {
		args := c.PostArgs()

		for i := range dec {
			v := fast.BytesToString(args.Peek(dec[i].name))

			if err = dec[i].scan(unsafe.Add(p, dec[i].offset), v); err != nil {
				return
			}
		}

		return
	}

	return
}
