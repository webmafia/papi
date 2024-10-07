package openapi

import jsoniter "github.com/json-iterator/go"

type SchemaType string

const (
	Array   SchemaType = "array"
	Boolean SchemaType = "boolean"
	Integer SchemaType = "integer"
	Number  SchemaType = "number"
	Object  SchemaType = "object"
	String  SchemaType = "string"
)

func (p SchemaType) String() string {
	return string(p)
}

func (p SchemaType) JsonEncode(_ *encoderContext, s *jsoniter.Stream) {
	s.WriteString(p.String())
}
