package registry

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/webbmaffian/papi/openapi"
)

func (r *Registry) Schema(typ reflect.Type, tag ...reflect.StructTag) (schema openapi.Schema, err error) {
	var tags reflect.StructTag

	if len(tag) > 0 {
		tags = tag[0]
	}

	schema, ok := r.getSchema(typ, tags)

	if !ok {
		schema, err = r.createSchema(typ, tags)
	}

	return
}

func (r *Registry) getSchema(typ reflect.Type, tags reflect.StructTag) (schema openapi.Schema, ok bool) {
	val, ok := r.typ[typ]

	if ok {
		schema = val.Describe(tags)
	}

	return
}

func (r *Registry) createSchema(typ reflect.Type, tags reflect.StructTag) (openapi.Schema, error) {
	switch kind := typ.Kind(); kind {

	case reflect.Bool:
		return &openapi.Boolean{}, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &openapi.Integer{}, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &openapi.Integer{}, nil

	case reflect.Float32, reflect.Float64:
		return &openapi.Number{}, nil

	case reflect.Array:
		itemType, err := r.Schema(typ.Elem(), tags)

		if err != nil {
			return nil, err
		}

		return &openapi.Array{
			Items: itemType,
			Min:   typ.Len(),
			Max:   typ.Len(),
		}, nil

	case reflect.Pointer:
		return r.Schema(typ.Elem(), tags)

	case reflect.Slice:
		itemType, err := r.Schema(typ.Elem(), tags)

		if err != nil {
			return nil, err
		}

		return &openapi.Array{Items: itemType}, nil

	case reflect.String:
		return &openapi.String{}, nil

	case reflect.Struct:
		numFlds := typ.NumField()

		obj := &openapi.Object{
			Properties: make([]openapi.ObjectProperty, 0, numFlds),
		}

		for i := range numFlds {
			fld := typ.Field(i)

			if !fld.IsExported() {
				continue
			}

			name := fld.Name

			if jsonTag, ok := fld.Tag.Lookup("json"); ok {
				name, _, _ = strings.Cut(jsonTag, ",")
			}

			propSchema, err := r.Schema(fld.Type, fld.Tag)

			if err != nil {
				return nil, err
			}

			obj.Properties = append(obj.Properties, openapi.ObjectProperty{
				Name:   name,
				Schema: propSchema,
			})
		}

		if name := typ.Name(); name != "" {
			return &openapi.Ref{
				Name:   name,
				Schema: obj,
			}, nil
		}

		return obj, nil

	default:
		return nil, fmt.Errorf("cannot create schema for type: %s", kind.String())
	}
}

func (s *Registry) DescribeOperation(op *openapi.Operation, in, out reflect.Type) (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Input
	if creator, ok := s.req[in]; ok {
		err = creator.Describe(op, in)
	} else if s.def != nil {
		err = s.def.Describe(op, in)
	} else {
		err = errors.New("no input descriptor could be found nor created")
	}

	if err != nil {
		return
	}

	// Output
	if op.Response, err = s.Schema(out); err != nil {
		return
	}

	return
}
