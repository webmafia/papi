package request

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
	"github.com/webmafia/fastapi/internal"
	"github.com/webmafia/fastapi/scanner"
	"github.com/webmafia/fastapi/scanner/structs"
	"github.com/webmafia/fastapi/scanner/value"
)

type inputTags struct {
	Body  string `tag:"body" eq:"json"`
	Param string `tag:"param"`
	Query string `tag:"query"`
}
type fieldScanner struct {
	offset uintptr
	scan   scanner.RequestScanner
}

var _ scanner.RequestScannerCreator = (*requestScanner)(nil)

type requestScanner struct {
	reg     *scanner.Registry
	tagScan value.ValueScanner
}

func NewRequestScanner(r *scanner.Registry) (creator scanner.RequestScannerCreator, err error) {
	tagScan, err := structs.CreateTagScanner(r, internal.ReflectType[inputTags]())

	if err != nil {
		return
	}

	creator = &requestScanner{
		reg:     r,
		tagScan: tagScan,
	}

	return
}

// CreateScanner implements scanner.RequestScannerCreator.
func (r *requestScanner) CreateScanner(typ reflect.Type, tags reflect.StructTag, paramKeys []string) (scan scanner.RequestScanner, err error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("invalid struct")
	}

	numFields := typ.NumField()
	flds := make([]fieldScanner, 0, numFields)

	for i := 0; i < numFields; i++ {
		var sc scanner.RequestScanner
		var tags inputTags

		fld := typ.Field(i)

		if sc, ok := api.scanners.get(fld.Type); ok {
			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
			continue
		}

		if err = r.tagScan(unsafe.Pointer(&tags), string(fld.Tag)); err != nil {
			return
		}

		if tags.Body == "json" {
			if sc, err = createJsonScanner(api, fld.Type); err != nil {
				return
			}

			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		}

		if tags.Param != "" {
			idx := slices.Index(params, tags.Param)

			if idx < 0 {
				err = fmt.Errorf("unknown param '%s'", tags.Param)
				return
			}

			if sc, err = createParamScanner(api, fld.Type, idx); err != nil {
				return
			}

			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		}

		if tags.Query != "" {
			if sc, err = createQueryScanner(api, fld.Type, tags.Query); err != nil {
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

func createJsonScanner(api *API, typ reflect.Type) (scan RequestScanner, err error) {
	dec := api.opt.JsonPool.DecoderOf(typ)
	scan = func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		iter := api.opt.JsonPool.AcquireIterator(c.Request.BodyStream())
		defer api.opt.JsonPool.ReleaseIterator(iter)

		dec.Decode(p, iter)
		return iter.Error
	}

	return
}

func createParamScanner(api *API, typ reflect.Type, idx int) (scan RequestScanner, err error) {
	sc, err := api.opt.StringScan.Scanner(typ)

	if err != nil {
		return
	}

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		params := RequestParams(c)
		return sc(p, params.Value(idx))
	}, nil
}

func createQueryScanner(api *API, typ reflect.Type, key string) (scan RequestScanner, err error) {
	sc, err := api.opt.StringScan.Scanner(typ)

	if err != nil {
		return
	}

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		val := c.QueryArgs().Peek(key)

		if len(val) > 0 {
			return sc(p, fast.BytesToString(val))
		}

		return nil
	}, nil
}