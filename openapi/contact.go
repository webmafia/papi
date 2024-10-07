package openapi

import jsoniter "github.com/json-iterator/go"

type Contact struct {
	Name  string
	Url   string
	Email string
}

func (c *Contact) JsonEncode(_ *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("name")
	s.WriteString(c.Name)

	s.WriteMore()
	s.WriteObjectField("url")
	s.WriteString(c.Url)

	s.WriteMore()
	s.WriteObjectField("email")
	s.WriteString(c.Email)

	s.WriteObjectEnd()
}
