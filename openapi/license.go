package openapi

import jsoniter "github.com/json-iterator/go"

type License struct {
	Name       string
	Identifier string
	Url        string
}

func (l *License) JsonEncode(_ *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("name")
	s.WriteString(l.Name)

	if l.Url != "" {
		s.WriteMore()
		s.WriteObjectField("url")
		s.WriteString(l.Url)
	} else if l.Identifier != "" {
		s.WriteMore()
		s.WriteObjectField("identifier")
		s.WriteString(l.Identifier)
	}

	s.WriteObjectEnd()
}
