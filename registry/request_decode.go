package registry

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/internal"
)

type RequestDecoder func(p unsafe.Pointer, c *fasthttp.RequestCtx) error

type inputTags struct {
	Body     string `tag:"body" enum:"json"`
	Param    string `tag:"param"`
	Query    string `tag:"query"`
	Security string `tag:"security"`
}

func (t inputTags) IsZero() bool {
	return t.Body == "" && t.Param == "" && t.Query == "" && t.Security == ""
}

type fieldScanner struct {
	offset uintptr
	scan   RequestDecoder
}

func (r *Registry) CreateRequestDecoder(typ reflect.Type, paramKeys []string, caller *runtime.Func) (scan RequestDecoder, err error) {
	return r.createRequestDecoder(typ, paramKeys, caller)
}

func (r *Registry) createRequestDecoder(typ reflect.Type, paramKeys []string, caller *runtime.Func) (scan RequestDecoder, err error) {
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

		// If there are no tags, and the field type is a struct, dive into it
		if tags.IsZero() {
			if fld.Type.Kind() == reflect.Struct {
				fldScan, err := r.createRequestDecoder(fld.Type, paramKeys, caller)

				if err != nil {
					return nil, err
				}

				flds = append(flds, fieldScanner{
					offset: fld.Offset,
					scan:   fldScan,
				})
			}

			continue
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

			if sc, err = r.createParamDecoder(fld.Type, tags.Param, idx, fld.Tag); err != nil {
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

		if tags.Security != "" {
			action, resource, ok := strings.Cut(tags.Security, ":")

			if !ok {
				resource = strings.ToLower(internal.CallerTypeFromFunc(caller))
			}

			if err = r.policies.Register(action, resource, fld.Type); err != nil {
				return
			}

			if sc, err = r.createSecurityDecoder(fld.Type, action, resource); err != nil {
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
