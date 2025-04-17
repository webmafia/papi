package registry

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/webmafia/papi/internal"
	"github.com/webmafia/papi/internal/iterate"
	"github.com/webmafia/papi/openapi"
)

func (r *Registry) Schema(typ reflect.Type, tag ...reflect.StructTag) (schema openapi.Schema, err error) {
	var tags reflect.StructTag

	if len(tag) > 0 {
		tags = tag[0]
	}

	// 1. If there is an explicit registered schema describer, use it
	if schema, err, ok := r.getSchema(typ, tags); ok {
		return schema, err
	}

	// 2. If the type can describe itself, let it
	if typ.Implements(typeDescriber) {
		if v, ok := reflect.New(typ).Interface().(TypeDescriber); ok {
			if desc := v.TypeDescription(r); desc.Schema != nil {
				return desc.Schema(tags)
			}
		}
	}

	// 3. In all other cases, describe the schema automatically
	return r.describeSchema(typ, tags)
}

func (r *Registry) getSchema(typ reflect.Type, tags reflect.StructTag) (schema openapi.Schema, err error, ok bool) {
	desc, ok := r.desc[typ]

	if ok && desc.Schema != nil {
		schema, err = desc.Schema(tags)
	}

	return
}

func (r *Registry) describeSchema(typ reflect.Type, tags reflect.StructTag) (openapi.Schema, error) {
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
		fldName := tags.Get("body")

		if fldName == "" {
			fldName = "json"
		}

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

			if jsonTag, ok := fld.Tag.Lookup(fldName); ok {
				name, _, _ = strings.Cut(jsonTag, ",")

				if name == "-" {
					continue
				}
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
				if iterate.HasFlag(flags, "required") {
					obj.Required = append(obj.Required, name)
				}
			}

		}

		if internal.IsPublicType(typ) {
			return &openapi.Ref{
				Name:   typ.Name(),
				Schema: obj,
			}, nil
		}

		return scanSchemaTags(r, obj, tags)

	default:
		return nil, fmt.Errorf("cannot create schema for type: %s", kind.String())
	}
}
