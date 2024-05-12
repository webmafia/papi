package fastapi

import (
	"errors"
	"reflect"

	"github.com/gosimple/slug"
	"github.com/webmafia/fastapi/internal"
	"github.com/webmafia/fastapi/spec"
)

func addRouteDocs[U, I, O any](api *API[U], r Route[U, I, O]) (err error) {
	if api.docs == nil {
		return
	}

	op := spec.Operation{
		Id:          slug.Make(r.Summary),
		Path:        r.Path,
		Method:      string(r.Method),
		Summary:     r.Summary,
		Description: r.Description,
		Tags:        r.Tags,
	}

	var in I

	typ := reflect.TypeOf(in)

	if typ.Kind() != reflect.Struct {
		return errors.New("expected struct")
	}

	numFld := typ.NumField()

	for i := 0; i < numFld; i++ {
		fld := typ.Field(i)

		if fld.Name == "Body" {
			spec.SchemaFromStruct(fld.Type, api.docs.Schemas)
		}

		internal.IterateStructTags(fld.Tag, func(key, val string) bool {
			switch key {

			case "param":
				op.Parameters = append(op.Parameters, spec.Parameter{
					Name:        val,
					In:          spec.InPath,
					Description: fld.Tag.Get("docs"),
					Required:    true,
				})

			case "query":
				op.Parameters = append(op.Parameters, spec.Parameter{
					Name:        val,
					In:          spec.InQuery,
					Description: fld.Tag.Get("docs"),
					// TODO: Required
				})

			}
			return true
		})
	}

	api.docs.Paths = append(api.docs.Paths, op)

	return
}
