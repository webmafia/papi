package openapi

import jsoniter "github.com/json-iterator/go"

type ParameterIn string

const (
	InQuery  ParameterIn = "query"
	InHeader ParameterIn = "header"
	InPath   ParameterIn = "path"
	InCookie ParameterIn = "cookie"
)

func (p ParameterIn) String() string {
	return string(p)
}

func (p ParameterIn) JsonEncode(_ *encoderContext, s *jsoniter.Stream) {
	s.WriteString(p.String())
}
