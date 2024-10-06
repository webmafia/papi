package request

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fastapi/internal"
	"github.com/webmafia/fastapi/pool/json"
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
	json    *json.Pool
	tagScan value.ValueScanner
}

func NewRequestScanner(r *scanner.Registry, json *json.Pool) (creator scanner.RequestScannerCreator, err error) {
	tagScan, err := structs.CreateTagScanner(r, internal.ReflectType[inputTags]())

	if err != nil {
		return
	}

	creator = &requestScanner{
		reg:     r,
		json:    json,
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

		if sc, err := r.reg.CreateRequestScanner(fld.Type, fld.Tag, paramKeys); err == nil {
			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		}

		if err = r.tagScan(unsafe.Pointer(&tags), string(fld.Tag)); err != nil {
			return
		}

		if tags.Body == "json" {
			if sc, err = r.createJsonScanner(fld.Type); err != nil {
				return
			}

			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		}

		if tags.Param != "" {
			idx := slices.Index(paramKeys, tags.Param)

			if idx < 0 {
				err = fmt.Errorf("unknown param '%s'", tags.Param)
				return
			}

			if sc, err = r.createParamScanner(fld.Type, idx); err != nil {
				return
			}

			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		}

		if tags.Query != "" {
			if sc, err = r.createQueryScanner(fld.Type, tags.Query); err != nil {
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
