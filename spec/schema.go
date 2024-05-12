package spec

import (
	"errors"
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

/*
"description": "",
"enum": [""],
"format": "",
"items": {},
"maximum": 0,
"maxItems": 0,
"maxLength": 0,
"minimum": 0,
"minItems": 0,
"minLength": 0,
"nullable": false,
"pattern": "",
"properties": {},
"readOnly": false,
"required": [],
"type": "",
"uniqueItems": false,
"writeOnly": false,
*/

type Schema struct {
	Title       string
	Description string
	Type        SchemaType
	Required    []string
	Enum        []string
	Format      string // string format
	Pattern     string // Regex
	Min         int
	Max         int
	Items       *Schema
	Properties  map[string]*Schema
	Nullable    bool
	ReadOnly    bool
	WriteOnly   bool
	UniqueItems bool
}

func SchemaFromStruct(typ reflect.Type, schemas map[reflect.Type]*Schema) (sch *Schema, err error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("expected struct")
	}

	if sch, ok := schemas[typ]; ok {
		return sch, nil
	}

	sch = &Schema{
		Title: typ.Name(),
		Type:  Object,
	}

	numFld := typ.NumField()

	for i := 0; i < numFld; i++ {
		fld := typ.Field(i)

	}

	return
}

func (sch *Schema) JsonEncode(ctx *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	sch.Type.JsonEncode(ctx, s)

	if sch.Title != "" {
		s.WriteMore()
		s.WriteObjectField("title")
		s.WriteString(sch.Title)
	}

	if sch.Description != "" {
		s.WriteMore()
		s.WriteObjectField("description")
		s.WriteString(sch.Description)
	}

	if sch.Nullable {
		s.WriteMore()
		s.WriteObjectField("nullable")
		s.WriteBool(sch.Nullable)
	}

	if sch.ReadOnly {
		s.WriteMore()
		s.WriteObjectField("readOnly")
		s.WriteBool(sch.ReadOnly)
	}

	if sch.WriteOnly {
		s.WriteMore()
		s.WriteObjectField("writeOnly")
		s.WriteBool(sch.WriteOnly)
	}

	switch sch.Type {

	case Array:
		if sch.Min > 0 {
			s.WriteMore()
			s.WriteObjectField("minItems")
			s.WriteInt(sch.Min)
		}

		if sch.Max > 0 {
			s.WriteMore()
			s.WriteObjectField("maxItems")
			s.WriteInt(sch.Max)
		}

		if sch.Items != nil {
			s.WriteMore()
			s.WriteObjectField("items")
			sch.Items.JsonEncode(ctx, s)
		}

		if sch.UniqueItems {
			s.WriteMore()
			s.WriteObjectField("uniqueItems")
			s.WriteBool(sch.UniqueItems)
		}

	case Integer, Number:
		if sch.Min >= 0 {
			s.WriteMore()
			s.WriteObjectField("minimum")
			s.WriteInt(sch.Min)
		}

		if sch.Max > 0 {
			s.WriteMore()
			s.WriteObjectField("maximum")
			s.WriteInt(sch.Max)
		}

	case Object:
		if len(sch.Required) > 0 {
			s.WriteMore()
			s.WriteObjectField("required")
			s.WriteArrayStart()

			for i := range sch.Required {
				if i != 0 {
					s.WriteMore()
				}

				s.WriteString(sch.Required[i])
			}

			s.WriteArrayEnd()
		}

		if len(sch.Properties) > 0 {
			s.WriteMore()
			s.WriteObjectField("properties")
			s.WriteObjectStart()

			var written bool

			for k, v := range sch.Properties {
				if written {
					s.WriteMore()
				} else {
					written = true
				}

				s.WriteObjectField(k)
				v.JsonEncode(ctx, s)
			}

			s.WriteObjectEnd()
		}

	case String:
		if len(sch.Enum) > 0 {
			s.WriteMore()
			s.WriteObjectField("enum")
			s.WriteArrayStart()

			for i := range sch.Enum {
				if i != 0 {
					s.WriteMore()
				}

				s.WriteString(sch.Enum[i])
			}

			s.WriteArrayEnd()
		}

		if sch.Format != "" {
			s.WriteMore()
			s.WriteObjectField("format")
			s.WriteString(sch.Format)
		}

		if sch.Pattern != "" {
			s.WriteMore()
			s.WriteObjectField("pattern")
			s.WriteString(sch.Pattern)
		}

		if sch.Min > 0 {
			s.WriteMore()
			s.WriteObjectField("minLength")
			s.WriteInt(sch.Min)
		}

		if sch.Max > 0 {
			s.WriteMore()
			s.WriteObjectField("maxLength")
			s.WriteInt(sch.Max)
		}

	}

	s.WriteObjectEnd()
}
