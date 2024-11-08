package openapi

import jsoniter "github.com/json-iterator/go"

type Parameter struct {
	Name        string
	In          ParameterIn
	Description string
	Schema      Schema
	Required    bool
}

func (p *Parameter) JsonEncode(ctx *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("name")
	s.WriteString(p.Name)

	s.WriteMore()
	s.WriteObjectField("in")
	p.In.JsonEncode(ctx, s)

	s.WriteMore()
	s.WriteObjectField("description")
	s.WriteString(p.Description)

	s.WriteMore()
	s.WriteObjectField("required")
	s.WriteBool(p.Required)

	if _, ok := p.Schema.(*Array); ok {
		s.WriteMore()
		s.WriteObjectField("explode")
		s.WriteBool(false)
	}

	s.WriteMore()
	s.WriteObjectField("schema")
	p.Schema.encodeSchema(ctx, s)

	s.WriteObjectEnd()
}
