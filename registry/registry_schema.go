package registry

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/webbmaffian/papi/internal"
	"github.com/webbmaffian/papi/openapi"
)

type ParamSchemer interface {
	ParamSchema(reg *Registry, tags reflect.StructTag) (schema openapi.Schema, err error)
}

var paramSchemer = reflect.TypeFor[ParamSchemer]()

func (r *Registry) Schema(typ reflect.Type, tag ...reflect.StructTag) (schema openapi.Schema, err error) {
	var tags reflect.StructTag

	if len(tag) > 0 {
		tags = tag[0]
	}

	if schema, ok := r.getSchema(typ, tags); ok {
		return schema, nil
	}

	if typ.Implements(paramSchemer) {
		if schemer, ok := reflect.New(typ).Interface().(ParamSchemer); ok {
			return schemer.ParamSchema(r, tags)
		}
	}

	return r.createSchema(typ, tags)
}

func (r *Registry) getSchema(typ reflect.Type, tags reflect.StructTag) (schema openapi.Schema, ok bool) {
	val, ok := r.typ[typ]

	if ok {
		schema = val.ParamSchema(tags)
	}

	return
}

func (r *Registry) createSchema(typ reflect.Type, tags reflect.StructTag) (openapi.Schema, error) {
	switch kind := typ.Kind(); kind {

	case reflect.Bool:
		return scanSchemaTags(r, &openapi.Boolean{}, tags)

	case reflect.Int:
		return scanSchemaTags(r, &openapi.Integer[int]{
			Min: math.MinInt,
			Max: math.MaxInt,
		}, tags)

	case reflect.Int8:
		return scanSchemaTags(r, &openapi.Integer[int8]{
			Min: math.MinInt8,
			Max: math.MaxInt8,
		}, tags)

	case reflect.Int16:
		return scanSchemaTags(r, &openapi.Integer[int16]{
			Min: math.MinInt16,
			Max: math.MaxInt16,
		}, tags)

	case reflect.Int32:
		return scanSchemaTags(r, &openapi.Integer[int32]{
			Min: math.MinInt32,
			Max: math.MaxInt32,
		}, tags)

	case reflect.Int64:
		return scanSchemaTags(r, &openapi.Integer[int64]{
			Min: math.MinInt64,
			Max: math.MaxInt64,
		}, tags)

	case reflect.Uint:
		return scanSchemaTags(r, &openapi.Integer[uint]{
			Max: math.MaxUint,
		}, tags)

	case reflect.Uint8:
		return scanSchemaTags(r, &openapi.Integer[uint8]{
			Max: math.MaxUint8,
		}, tags)

	case reflect.Uint16:
		return scanSchemaTags(r, &openapi.Integer[uint16]{
			Max: math.MaxUint16,
		}, tags)

	case reflect.Uint32:
		return scanSchemaTags(r, &openapi.Integer[uint32]{
			Max: math.MaxUint32,
		}, tags)

	case reflect.Uint64:
		return scanSchemaTags(r, &openapi.Integer[uint64]{
			Max: math.MaxUint64,
		}, tags)

	case reflect.Float32:
		return scanSchemaTags(r, &openapi.Number[float32]{
			Min: -math.MaxFloat32,
			Max: math.MaxFloat32,
		}, tags)

	case reflect.Float64:
		return scanSchemaTags(r, &openapi.Number[float64]{
			Min: -math.MaxFloat64,
			Max: math.MaxFloat64,
		}, tags)

	case reflect.Array:
		itemType, err := r.Schema(typ.Elem(), tags)

		if err != nil {
			return nil, err
		}

		return scanSchemaTags(r, &openapi.Array{
			Items: itemType,
			Min:   typ.Len(),
			Max:   typ.Len(),
		}, tags)

	case reflect.Pointer:
		return r.Schema(typ.Elem(), tags)

	case reflect.Slice:
		itemType, err := r.Schema(typ.Elem(), tags)

		if err != nil {
			return nil, err
		}

		return scanSchemaTags(r, &openapi.Array{Items: itemType}, tags)

	case reflect.String:
		return scanSchemaTags(r, &openapi.String{}, tags)

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

			if flags, ok := fld.Tag.Lookup("flags"); ok {
				if internal.HasFlag(flags, "required") {
					obj.Required = append(obj.Required, name)
				}
			}

		}

		if name := typ.Name(); name != "" {
			return &openapi.Ref{
				Name:   name,
				Schema: obj,
			}, nil
		}

		return scanSchemaTags(r, obj, tags)

	default:
		return nil, fmt.Errorf("cannot create schema for type: %s", kind.String())
	}
}

func (s *Registry) DescribeOperation(op *openapi.Operation, in, out reflect.Type) (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Input
	if creator, ok := s.req[in]; ok {
		err = creator.DescribeOperation(op)
	} else if s.def != nil {
		err = s.def.DescribeOperation(op, in)
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
