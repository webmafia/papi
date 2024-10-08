package registry

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/webmafia/fastapi/openapi"
)

func (r *Registry) Schema(typ reflect.Type) (schema *openapi.Schema, err error) {
	schema, ok := r.getSchema(typ)

	if !ok {
		schema, err = r.createSchema(typ)
	}

	return
}

func (r *Registry) getSchema(typ reflect.Type) (schema *openapi.Schema, ok bool) {
	// r.mu.RLock()
	// defer r.mu.RUnlock()

	val, ok := r.typ[typ]

	if ok {
		schema = val.Describe("")
	}

	return
}

func (r *Registry) createSchema(typ reflect.Type) (s *openapi.Schema, err error) {
	// r.mu.Lock()
	// defer r.mu.Unlock()

	s = new(openapi.Schema)

	if err = r.describeSchema(s, typ); err != nil {
		return
	}

	return
}

func (r *Registry) describeSchema(s *openapi.Schema, typ reflect.Type) (err error) {
	switch kind := typ.Kind(); kind {

	case reflect.Bool:
		s.Type = openapi.Boolean

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s.Type = openapi.Integer

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s.Type = openapi.Integer

	case reflect.Float32, reflect.Float64:
		s.Type = openapi.Number

	case reflect.Array:
		s.Type = openapi.Array
		s.Min = typ.Len()
		s.Max = s.Min
		s.Items = new(openapi.Schema)

		if err = r.describeSchema(s.Items, typ.Elem()); err != nil {
			return
		}

	case reflect.Pointer:
		return r.describeSchema(s, typ.Elem())

	case reflect.Slice:
		s.Type = openapi.Array
		s.Items = new(openapi.Schema)

		if err = r.describeSchema(s.Items, typ.Elem()); err != nil {
			return
		}

	case reflect.String:
		s.Type = openapi.String

	case reflect.Struct:
		s.Title = typ.Name()
		numFlds := typ.NumField()
		s.Type = openapi.Object
		s.Properties = make([]openapi.Property, 0, numFlds)
		s.ShouldBeRef = true

		for i := range numFlds {
			fld := typ.Field(i)

			if !fld.IsExported() {
				continue
			}

			name := fld.Name

			if jsonTag, ok := fld.Tag.Lookup("json"); ok {
				name, _, _ = strings.Cut(jsonTag, ",")
			}

			prop := openapi.Property{
				Name: name,
			}

			if prop.Schema, err = r.Schema(fld.Type); err != nil {
				return
			}

			s.Properties = append(s.Properties, prop)

		}

	default:
		return fmt.Errorf("cannot create schema for type: %s", kind.String())
	}

	return
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
