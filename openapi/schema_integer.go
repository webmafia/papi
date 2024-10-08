package openapi

import jsoniter "github.com/json-iterator/go"

var _ Schema = (*Integer)(nil)

type Integer struct {
	Title       string
	Description string
	Min         int
	Max         int
	Nullable    bool
	ReadOnly    bool
	WriteOnly   bool
}

func (sch *Integer) encodeSchema(ctx *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("integer")

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

	s.WriteObjectEnd()
}
