package papi

import jsoniter "github.com/json-iterator/go"

type JsonEncoder interface {
	EncodeJson(s *jsoniter.Stream) error
}
