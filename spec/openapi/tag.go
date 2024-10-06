package openapi

import jsoniter "github.com/json-iterator/go"

type Tag struct {
	Name        string
	Description string
}

func NewTag(name, description string) *Tag {
	return &Tag{
		Name:        name,
		Description: description,
	}
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
