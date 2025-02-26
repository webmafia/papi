package registry

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"slices"
	"unsafe"

	"github.com/valyala/fasthttp"
)

type RequestDecoder func(p unsafe.Pointer, c *fasthttp.RequestCtx) error

type inputTags struct {
	Body       string `tag:"body" enum:"json"`
	Param      string `tag:"param"`
	Query      string `tag:"query"`
	Permission string `tag:"perm"`
}

func (t inputTags) IsZero() bool {
	return t.Body == "" && t.Param == "" && t.Query == "" && t.Permission == ""
}

type fieldScanner struct {
	offset uintptr
	scan   RequestDecoder
}

func (r *Registry) CreateRequestDecoder(typ reflect.Type, paramKeys []string, caller *runtime.Func) (scan RequestDecoder, perm string, err error) {
	return r.createRequestDecoder(typ, paramKeys, caller)
}

func (r *Registry) createRequestDecoder(typ reflect.Type, paramKeys []string, caller *runtime.Func) (scan RequestDecoder, perm string, err error) {
	if typ.Kind() != reflect.Struct {
		err = errors.New("invalid struct")
		return
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
				fldScan, subPerm, err := r.createRequestDecoder(fld.Type, paramKeys, caller)

				if err != nil {
					return nil, "", err
				}

				if subPerm != "" {
					if perm != "" {
						return nil, "", fmt.Errorf("%s.%s in %s has a permission tag, but permission for the route is already set", typ.Name(), fld.Name, caller.Name())
					}

					perm = subPerm
				}

				flds = append(flds, fieldScanner{
					offset: fld.Offset,
					scan:   fldScan,
				})
			}

			continue
		}

		if tags.Body != "" {
			if sc, err = r.getCustomDecoder(fld.Type, fld.Tag); err != nil {
				return
			}

			if sc == nil {
				if tags.Body != "json" {
					return nil, "", fmt.Errorf("unknown body type: '%s'", tags.Body)
				}

				if sc, err = r.createJsonDecoder(fld.Type); err != nil {
					return
				}
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

		if tags.Permission != "" && r.securityScheme != nil {
			if sc, perm, err = r.securityScheme.OperationSecurityHandler(fld.Type, tags.Permission, caller); err != nil {
				return
			}

			if sc != nil {
				flds = append(flds, fieldScanner{
					offset: fld.Offset,
					scan:   sc,
				})
			}
		}
	}

	if r.forcePermTag && perm == "" && r.securityScheme != nil {
		return nil, "", fmt.Errorf("route %s is missing a permission tag, which is required when a security scheme is set", caller.Name())
	}

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) (err error) {
		for _, fld := range flds {
			if err = fld.scan(unsafe.Add(p, fld.offset), c); err != nil {
				return
			}
		}

		return
	}, perm, nil
}
