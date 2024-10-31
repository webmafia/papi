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
	"github.com/webmafia/papi/security"
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

func (r *Registry) CreateRequestDecoder(typ reflect.Type, paramKeys []string, caller *runtime.Func) (scan RequestDecoder, perm security.Permission, err error) {
	return r.createRequestDecoder(typ, paramKeys, caller)
}

func (r *Registry) createRequestDecoder(typ reflect.Type, paramKeys []string, caller *runtime.Func) (scan RequestDecoder, perm security.Permission, err error) {
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

		if tags.Body == "json" {
			if sc, err = r.createJsonDecoder(fld.Type); err != nil {
				return
			}

			flds = append(flds, fieldScanner{
				offset: fld.Offset,
				scan:   sc,
			})
		} else if tags.Body != "" {
			return nil, "", errors.New("the only valid 'body' tag value is 'json'")
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

		if tags.Permission != "" {
			if r.forcePermTag && r.gatekeeper == nil {
				return nil, "", fmt.Errorf("%s.%s in %s has permission tag '%s', but no API gatekeeper is set", typ.Name(), fld.Name, caller.Name(), tags.Permission)
			}

			if tags.Permission != "-" {
				perm = security.Permission(tags.Permission)

				if !perm.HasResource() {
					perm.SetResource(strings.ToLower(internal.CallerTypeFromFunc(caller)))
				}

				if err = r.gatekeeper.RegisterPermission(perm, fld.Type); err != nil {
					return
				}

				if r.gatekeeper != nil {
					if sc, err = r.createPermissionDecoder(fld.Type, perm); err != nil {
						return
					}

					flds = append(flds, fieldScanner{
						offset: fld.Offset,
						scan:   sc,
					})
				}
			}
		}
	}

	if r.forcePermTag && perm != "" && r.gatekeeper != nil {
		return nil, "", fmt.Errorf("route %s is missing a permission tag, which is required when an API gatekeeper is set", caller.Name())
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
