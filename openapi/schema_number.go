package openapi

import jsoniter "github.com/json-iterator/go"

var _ Schema = (*Number)(nil)

type Number struct {
	Title       string
	Description string
	Min         float64
	Max         float64
	Nullable    bool
	ReadOnly    bool
	WriteOnly   bool
}

func (sch *Number) encodeSchema(ctx *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("number")

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
		s.WriteFloat64(sch.Min)
	}

	if sch.Max > 0 {
		s.WriteMore()
		s.WriteObjectField("maximum")
		s.WriteFloat64(sch.Max)
	}

	s.WriteObjectEnd()
}
