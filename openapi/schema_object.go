package openapi

import jsoniter "github.com/json-iterator/go"

var _ Schema = (*Object)(nil)

type Object struct {
	Title       string
	Description string
	Required    []string
	Properties  []ObjectProperty
	Nullable    bool
	ReadOnly    bool
	WriteOnly   bool
}

type ObjectProperty struct {
	Name   string
	Schema Schema
}

func (sch *Object) encodeSchema(ctx *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("object")

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
			sch.Properties[i].Schema.encodeSchema(ctx, s)
		}

		s.WriteObjectEnd()
	}

	s.WriteObjectEnd()
}
