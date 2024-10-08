package openapi

import jsoniter "github.com/json-iterator/go"

type Schema interface {
	encodeSchema(ctx *encoderContext, s *jsoniter.Stream)
}
