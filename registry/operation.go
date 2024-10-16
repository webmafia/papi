package registry

import (
	"errors"
	"reflect"

	"github.com/webbmaffian/papi/openapi"
)

func (s *Registry) DescribeOperation(op *openapi.Operation, in, out reflect.Type) (err error) {

	// Input
	if err = s.describeOperation(op, in); err != nil {
		return
	}

	// Output
	if op.Response, err = s.Schema(out); err != nil {
		return
	}

	return
}

func (r *Registry) describeOperation(op *openapi.Operation, typ reflect.Type) (err error) {
	if typ.Kind() != reflect.Struct {
		return errors.New("invalid struct")
	}

	numFields := typ.NumField()

	for i := 0; i < numFields; i++ {
		var tags inputTags
		fld := typ.Field(i)

		if err = ScanTags(r, &tags, fld.Tag); err != nil {
			return
		}

		if tags.Body == "json" {
			schema, err := r.Schema(fld.Type, fld.Tag)

			if err != nil {
				return err
			}

			op.RequestBody = schema
		}

		if tags.Param != "" {
			schema, err := r.Schema(fld.Type, fld.Tag)

			if err != nil {
				return err
			}

			param := openapi.Parameter{
				Name:     tags.Param,
				In:       openapi.InPath,
				Required: true,
				Schema:   schema,
			}

			op.Parameters = append(op.Parameters, param)
		}

		if tags.Query != "" {
			schema, err := r.Schema(fld.Type, fld.Tag)

			if err != nil {
				return err
			}

			param := openapi.Parameter{
				Name:   tags.Query,
				In:     openapi.InQuery,
				Schema: schema,
			}

			op.Parameters = append(op.Parameters, param)
		}
	}

	return
}
