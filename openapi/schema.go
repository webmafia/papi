package openapi

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/webbmaffian/papi/internal/hasher"
)

type Schema interface {
	hasher.Hashable
	encodeSchema(ctx *encoderContext, s *jsoniter.Stream) error
}
