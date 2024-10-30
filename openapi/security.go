package openapi

import (
	jsoniter "github.com/json-iterator/go"
)

type Security struct {
	Action   string
	Resource string
	Valid    bool
}

func (sec Security) IsZero() bool {
	return sec.Action == "" || sec.Resource == ""
}

func (sec Security) String() string {
	if sec.IsZero() {
		return ""
	}

	return sec.Action + ":" + sec.Resource
}

func (ss Security) JsonEncode(_ *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("token")
	s.WriteArrayStart()
	s.WriteString(ss.String())
	s.WriteArrayEnd()

	s.WriteObjectEnd()
}
