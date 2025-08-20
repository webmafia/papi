package registry

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
)

type multipartValueDecoder struct {
	offset uintptr
	scan   Decoder
	name   string
}

// type multipartFileDecoder struct {
// 	offset uintptr
// 	name   string
// }

func (r *Registry) createMultipartHandler(typ reflect.Type) (scan Handler, err error) {
	if typ.Kind() != reflect.Struct {
		err = errors.New("invalid struct")
		return
	}

	numFields := typ.NumField()
	valDec := make([]multipartValueDecoder, 0, numFields)

	for i := 0; i < numFields; i++ {
		fld := typ.Field(i)
		name := fld.Name

		if v := fld.Tag.Get("multipart"); v != "" {
			name = v
		} else if v := fld.Tag.Get("form"); v != "" {
			name = v
		}

		if name == "-" {
			continue
		}

		sc, err := r.Decoder(fld.Type, fld.Tag)

		if err != nil {
			return nil, err
		}

		valDec = append(valDec, multipartValueDecoder{
			offset: fld.Offset,
			scan:   sc,
			name:   name,
		})
	}

	scan = func(c *fasthttp.RequestCtx, p unsafe.Pointer) (err error) {
		form, err := c.MultipartForm()

		if err != nil {
			return
		}

		for i := range valDec {
			v, ok := form.Value[valDec[i].name]

			if !ok || len(v) < 1 {
				continue
			}

			if err = valDec[i].scan(unsafe.Add(p, valDec[i].offset), v[0]); err != nil {
				return
			}
		}

		return
	}

	return
}
