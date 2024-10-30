package openapi

import (
	jsoniter "github.com/json-iterator/go"
)

type Security struct {
	Scope string
}

func (sec Security) IsZero() bool {
	return sec.Scope == ""
}

func (ss Security) JsonEncode(_ *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	if !ss.IsZero() {
		s.WriteObjectField("token")
		s.WriteArrayStart()
		s.WriteString(ss.Scope)
		s.WriteArrayEnd()
	}

	s.WriteObjectEnd()
}
