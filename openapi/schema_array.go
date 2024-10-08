package openapi

import jsoniter "github.com/json-iterator/go"

var _ Schema = (*Array)(nil)

type Array struct {
	Title       string
	Description string
	Min         int
	Max         int
	Items       Schema
	Nullable    bool
	ReadOnly    bool
	WriteOnly   bool
	UniqueItems bool
}

func (sch *Array) encodeSchema(ctx *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("array")

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
		sch.Items.encodeSchema(ctx, s)
	}

	if sch.UniqueItems {
		s.WriteMore()
		s.WriteObjectField("uniqueItems")
		s.WriteBool(sch.UniqueItems)
	}

	s.WriteObjectEnd()
}
