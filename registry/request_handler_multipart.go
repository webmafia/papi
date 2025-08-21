package registry

import (
	"errors"
	"reflect"
	"slices"
	"unsafe"

	"github.com/valyala/fasthttp"
)

type multipartParser struct {
	offset uintptr
	parse  Parser
	name   string
}

type multipartBinder struct {
	offset uintptr
	bind   Binder
}

func (r *Registry) createMultipartBinder(typ reflect.Type) (_ Binder, err error) {
	if typ.Kind() != reflect.Struct {
		err = errors.New("invalid struct")
		return
	}

	numFields := typ.NumField()
	binders := make([]multipartBinder, 0, numFields)
	parsers := make([]multipartParser, 0, numFields)

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

		bind, err := r.Binder(fld.Type, fld.Tag)

		if err != nil {
			return nil, err
		}

		if bind != nil {
			binders = append(binders, multipartBinder{
				offset: fld.Offset,
				bind:   bind,
			})
		} else {
			sc, err := r.Parser(fld.Type, fld.Tag)

			if err != nil {
				return nil, err
			}

			if sc != nil {
				parsers = append(parsers, multipartParser{
					offset: fld.Offset,
					parse:  sc,
					name:   name,
				})
			}
		}
	}

	parsers = slices.Clip(parsers)
	binders = slices.Clip(binders)

	return func(c *fasthttp.RequestCtx, p unsafe.Pointer) (err error) {
		form, err := c.MultipartForm()

		if err != nil {
			return
		}

		for i := range parsers {
			v, ok := form.Value[parsers[i].name]

			if !ok || len(v) < 1 {
				continue
			}

			if err = parsers[i].parse(unsafe.Add(p, parsers[i].offset), v[0]); err != nil {
				return
			}
		}

		for i := range binders {
			if err = binders[i].bind(c, unsafe.Add(p, binders[i].offset)); err != nil {
				return
			}
		}

		return
	}, nil
}
