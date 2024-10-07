package openapi

import jsoniter "github.com/json-iterator/go"

type License struct {
	Name string
	Url  string
}

func (l *License) JsonEncode(_ *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("name")
	s.WriteString(l.Name)

	s.WriteMore()
	s.WriteObjectField("url")
	s.WriteString(l.Url)

	s.WriteObjectEnd()
}
