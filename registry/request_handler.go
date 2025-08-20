package registry

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"unsafe"

	"github.com/modern-go/reflect2"
	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/internal"
	"github.com/webmafia/papi/security"
)

type inputTags struct {
	Body       string `tag:"body" enum:"json"`
	Param      string `tag:"param"`
	Query      string `tag:"query"`
	Permission string `tag:"perm"`
}

func (t inputTags) IsZero() bool {
	return t.Body == "" && t.Param == "" && t.Query == "" && t.Permission == ""
}

type fieldHandler struct {
	offset  uintptr
	handler Binder
}

func (r *Registry) CreateBinder(typ reflect.Type, paramKeys []string, caller *runtime.Func) (scan Binder, perm string, err error) {
	return r.createBinder(typ, paramKeys, caller)
}

func (r *Registry) createBinder(typ reflect.Type, paramKeys []string, caller *runtime.Func) (scan Binder, perm string, err error) {
	if typ.Kind() != reflect.Struct {
		err = errors.New("invalid struct")
		return
	}

	numFields := typ.NumField()
	flds := make([]fieldHandler, 0, numFields+1)

	if r.gatekeeper != nil {
		flds = append(flds, fieldHandler{
			handler: func(c *fasthttp.RequestCtx, _ unsafe.Pointer) error {
				return r.gatekeeper.PreRequest(c)
			},
		})
	}

	for i := 0; i < numFields; i++ {
		var sc Binder
		var tags inputTags

		fld := typ.Field(i)

		if err = ScanTags(r, &tags, fld.Tag); err != nil {
			return
		}

		// If there are no tags, and the field type is a struct, dive into it
		if tags.IsZero() {
			if fld.Type.Kind() == reflect.Struct {
				fldScan, subPerm, err := r.createBinder(fld.Type, paramKeys, caller)

				if err != nil {
					return nil, "", err
				}

				if subPerm != "" {
					if perm != "" {
						return nil, "", fmt.Errorf("%s.%s in %s has a permission tag, but permission for the route is already set", typ.Name(), fld.Name, caller.Name())
					}

					perm = subPerm
				}

				flds = append(flds, fieldHandler{
					offset:  fld.Offset,
					handler: fldScan,
				})
			}

			continue
		}

		if tags.Body != "" {
			if sc, err = r.getCustomBinder(fld.Type, fld.Tag); err != nil {
				return
			}

			if sc == nil {
				switch tags.Body {
				case "json":
					sc, err = r.createJsonBinder(fld.Type)
				case "form":
					sc, err = r.createFormBinder(fld.Type)
				case "multipart":
					sc, err = r.createMultipartBinder(fld.Type)
				default:
					err = fmt.Errorf("unknown body type: '%s'", tags.Body)
				}

				if err != nil {
					return
				}
			}

			flds = append(flds, fieldHandler{
				offset:  fld.Offset,
				handler: sc,
			})
		}

		if tags.Param != "" {
			idx := slices.Index(paramKeys, tags.Param)

			if idx < 0 {
				err = fmt.Errorf("unknown param '%s'", tags.Param)
				return
			}

			if sc, err = r.createParamBinder(fld.Type, tags.Param, idx, fld.Tag); err != nil {
				return
			}

			flds = append(flds, fieldHandler{
				offset:  fld.Offset,
				handler: sc,
			})
		}

		if tags.Query != "" {
			if sc, err = r.createQueryBinder(fld.Type, tags.Query, fld.Tag); err != nil {
				return
			}

			flds = append(flds, fieldHandler{
				offset:  fld.Offset,
				handler: sc,
			})
		}

		if tags.Permission != "" && r.gatekeeper != nil {
			switch gk := r.gatekeeper.(type) {

			case security.RolesGatekeeper:
				if sc, perm, err = r.createOperationSecurityBinder(fld.Type, tags.Permission, caller, gk); err != nil {
					return
				}

			case security.CustomGatekeeper:
				sc = func(c *fasthttp.RequestCtx, p unsafe.Pointer) error {
					return gk.HandleSecurity(c, security.Permission(tags.Permission), reflect.NewAt(fld.Type, p).Interface())
				}

			default:
				err = fmt.Errorf("unknown gatekeeper type: %T", gk)
				return

			}

			if sc != nil {
				flds = append(flds, fieldHandler{
					offset:  fld.Offset,
					handler: sc,
				})
			}
		}
	}

	if !r.OptionalPermTag() && perm == "" {
		return nil, "", fmt.Errorf("route %s is missing a permission tag, which is required by the Gatekeeper", caller.Name())
	}

	return func(c *fasthttp.RequestCtx, p unsafe.Pointer) (err error) {
		for _, fld := range flds {
			if err = fld.handler(c, unsafe.Add(p, fld.offset)); err != nil {
				return
			}
		}

		return
	}, perm, nil
}

func (r *Registry) createOperationSecurityBinder(typ reflect.Type, permTag string, caller *runtime.Func, gk security.RolesGatekeeper) (handler Binder, modTag string, err error) {
	if permTag == "-" {
		return nil, permTag, nil
	}

	perm := security.Permission(permTag)

	if !perm.HasResource() {
		perm.SetResource(strings.ToLower(internal.CallerTypeFromFunc(caller)))
	}

	if err = r.policies.Register(perm, typ); err != nil {
		return
	}

	typ2 := reflect2.Type2(typ)

	return func(c *fasthttp.RequestCtx, p unsafe.Pointer) (err error) {
		userRoles, err := gk.UserRoles(c)

		if err != nil {
			return
		}

		cond, err := r.policies.Get(userRoles, perm)

		if err != nil {
			return err
		}

		if cond != nil {
			typ2.UnsafeSet(p, cond)
		}

		return nil
	}, perm.String(), nil
}
