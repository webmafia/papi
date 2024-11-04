package openapi

import (
	jsoniter "github.com/json-iterator/go"
)

type SecurityRequirement struct {
	Name   string
	Scopes []string
}

func (sec SecurityRequirement) IsZero() bool {
	return sec.Name == ""
}

func (ss SecurityRequirement) JsonEncode(s *jsoniter.Stream) {
	s.WriteArrayStart()

	for i := range ss.Scopes {
		if i != 0 {
			s.WriteMore()
		}

		s.WriteString(ss.Scopes[i])
	}

	s.WriteArrayEnd()
}
