package spec

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/webmafia/fastapi/scan"
	"github.com/webmafia/fastapi/spec/schema"
)

var boolScan, _ = scan.CreateScanner(reflect.TypeOf(true))

func (d *Document) getSchema(typ reflect.Type) (sch schema.Schema, ok bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	sch, ok = d.Schemas[typ]
	return
}

func (d *Document) RegisterSchema(typ reflect.Type, sch schema.Schema) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.Schemas[typ] = sch
}

func (d *Document) SchemaOf(typ reflect.Type) (sch schema.Schema, err error) {
	if sch, ok := d.getSchema(typ); ok {
		return sch, nil
	}

	sch, err = d.generateSchema(typ)

	if err == nil {
		d.RegisterSchema(typ, sch)
	}

	return
}

func (d *Document) generateSchema(typ reflect.Type) (sch schema.Schema, err error) {
	switch kind := typ.Kind(); kind {

	case reflect.Bool:
		sch = &schema.Boolean{}

	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		sch = &schema.Integer{}

	case reflect.Float32,
		reflect.Float64:
		sch = &schema.Number{}

	case reflect.Pointer:
		return d.generateSchema(typ.Elem())

	case reflect.Slice,
		reflect.Array:
		elemSch, err := d.generateSchema(typ.Elem())

		if err != nil {
			return nil, err
		}

		sch = &schema.Array{
			Items: elemSch,
		}

	case reflect.String:
		sch = &schema.String{}

	case reflect.Struct:
		o := &schema.Object{}
		o.SetName(typ.Name())
		sch = o
		numFlds := typ.NumField()

		for i := 0; i < numFlds; i++ {
			fld := typ.Field(i)
			name := jsonFieldName(fld)

			if name == "" {
				continue
			}

			fldSch, err := d.generateSchema(fld.Type)

			if err != nil {
				return sch, err
			}

			if err = fldSch.ScanTags(fld.Tag); err != nil {
				return nil, err
			}

			if tag, ok := fld.Tag.Lookup("required"); ok {
				var req bool

				if err := boolScan(unsafe.Pointer(&req), tag); err != nil {
					return nil, err
				}

				if req {
					o.Required = append(o.Required, name)
				}
			}

			o.Properties = append(o.Properties, schema.ObjectProp{
				Name:   name,
				Schema: fldSch,
			})
		}

	default:
		return nil, fmt.Errorf("unsupported type: %s", kind)

	}

	return
}

func jsonFieldName(fld reflect.StructField) string {
	if !fld.IsExported() {
		return ""
	}

	if tag, ok := fld.Tag.Lookup("json"); ok {
		name, _, _ := strings.Cut(tag, ",")

		if name == "-" {
			return ""
		}

		if name != "" {
			return name
		}
	}

	return fld.Name
}
