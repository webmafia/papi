package spec

type ParameterIn string

const (
	InQuery  ParameterIn = "query"
	InHeader ParameterIn = "header"
	InPath   ParameterIn = "path"
	InCookie ParameterIn = "cookie"
)
