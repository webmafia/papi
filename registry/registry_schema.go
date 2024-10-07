package registry

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/webmafia/fastapi/openapi"
)

func (r *Registry) RegisterSchema(typ reflect.Type, schema *openapi.Schema) {
	// r.mu.Lock()
	// defer r.mu.Unlock()

	if schema == nil {
		delete(r.schemas, typ)
	} else {
		r.schemas[typ] = schema
	}
}

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

	schema, ok = r.schemas[typ]
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
				Name:   name,
				Schema: new(openapi.Schema),
			}

			if err = r.describeSchema(prop.Schema, fld.Type); err != nil {
				return
			}

			s.Properties = append(s.Properties, prop)

		}

	default:
		return fmt.Errorf("cannot create schema for type: %s", kind.String())
	}

	return
}
