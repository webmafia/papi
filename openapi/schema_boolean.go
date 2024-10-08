package openapi

import jsoniter "github.com/json-iterator/go"

var _ Schema = (*Boolean)(nil)

type Boolean struct {
	Title       string
	Description string
	Nullable    bool
	ReadOnly    bool
	WriteOnly   bool
}

func (sch *Boolean) encodeSchema(ctx *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("boolean")

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

	s.WriteObjectEnd()
}
