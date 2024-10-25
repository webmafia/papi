package openapi

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/papi/internal/hasher"
)

type Schema interface {
	hasher.Hashable
	GetTitle() string
	encodeSchema(ctx *encoderContext, s *jsoniter.Stream) error
	encodeValue(s *jsoniter.Stream, val string) error
}
