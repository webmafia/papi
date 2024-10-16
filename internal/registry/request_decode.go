package registry

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"unsafe"

	"github.com/valyala/fasthttp"
)

type RequestDecoder func(p unsafe.Pointer, c *fasthttp.RequestCtx) error

type inputTags struct {
	Body  string `tag:"body" enum:"json"`
	Param string `tag:"param"`
	Query string `tag:"query"`
}
type fieldScanner struct {
	offset uintptr
	scan   RequestDecoder
}

func (r *Registry) CreateRequestDecoder(typ reflect.Type, paramKeys []string) (scan RequestDecoder, err error) {
	return r.createRequestDecoder(typ, paramKeys)
}

func (r *Registry) createRequestDecoder(typ reflect.Type, paramKeys []string) (scan RequestDecoder, err error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("invalid struct")
	}

	numFields := typ.NumField()
	flds := make([]fieldScanner, 0, numFields)

	for i := 0; i < numFields; i++ {
		var sc RequestDecoder
		var tags inputTags

		fld := typ.Field(i)

		if err = ScanTags(r, &tags, fld.Tag); err != nil {
			return
		}

		if tags.Body == "json" {
			if sc, err = r.createJsonDecoder(fld.Type); err != nil {
				return
			}

			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		} else if tags.Body != "" {
			return nil, errors.New("the only valid 'body' tag value is 'json'")
		}

		if tags.Param != "" {
			idx := slices.Index(paramKeys, tags.Param)

			if idx < 0 {
				err = fmt.Errorf("unknown param '%s'", tags.Param)
				return
			}

			if sc, err = r.createParamDecoder(fld.Type, idx, fld.Tag); err != nil {
				return
			}

			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		}

		if tags.Query != "" {
			if sc, err = r.createQueryDecoder(fld.Type, tags.Query, fld.Tag); err != nil {
				return
			}

			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		}
	}

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) (err error) {
		for _, fld := range flds {
			if err = fld.scan(unsafe.Add(p, fld.offset), c); err != nil {
				return
			}
		}

		return
	}, nil
}
