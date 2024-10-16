package errors

import (
	jsoniter "github.com/json-iterator/go"
)

type ErrorDocumentor interface {
	Status() int
	ErrorDocument(s *jsoniter.Stream)
}
