package openapi

import jsoniter "github.com/json-iterator/go"

type Parameter struct {
	Name        string
	In          ParameterIn
	Description string
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

	s.WriteObjectEnd()
}
