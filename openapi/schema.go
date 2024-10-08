package openapi

import (
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
	Enum        []string `tag:"enum"`
	Format      string   `tag:"format"`
	Pattern     string   `tag:"pattern"`
	Min         int      `tag:"min"`
	Max         int      `tag:"max"`
	Items       *Schema
	Properties  []Property
	Nullable    bool `tag:"nullable"`
	ReadOnly    bool `tag:"readonly"`
	WriteOnly   bool `tag:"writeonly"`
	UniqueItems bool `tag:"unique"`
	ShouldBeRef bool
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

			for i := range sch.Properties {
				if i != 0 {
					s.WriteMore()
				}

				s.WriteObjectField(sch.Properties[i].Name)
				sch.Properties[i].Schema.JsonEncode(ctx, s)
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
