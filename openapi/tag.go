package openapi

import jsoniter "github.com/json-iterator/go"

type Tag struct {
	Name        string
	Description string
}

func NewTag(name string, description ...string) Tag {
	t := Tag{
		Name: name,
	}

	if len(description) > 0 {
		t.Description = description[0]
	}

	return t
}

func (t *Tag) JsonEncode(_ *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("name")
	s.WriteString(t.Name)

	if t.Description != "" {
		s.WriteMore()
		s.WriteObjectField("description")
		s.WriteString(t.Description)
	}

	s.WriteObjectEnd()
}
