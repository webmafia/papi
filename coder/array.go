package coder

import (
	jsoniter "github.com/json-iterator/go"
)

var _ Coder = (*Array)(nil)

type Array struct {
	General
	Min         int `tag:"min"`
	Max         int `tag:"max"`
	Items       Coder
	UniqueItems bool `tag:"unique"`
}

// EncodeSchema implements Coder.
func (a Array) EncodeSchema(s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("array")

	if a.Min > 0 {
		s.WriteMore()
		s.WriteObjectField("minItems")
		s.WriteInt(a.Min)
	}

	if a.Max > 0 {
		s.WriteMore()
		s.WriteObjectField("maxItems")
		s.WriteInt(a.Max)
	}

	if a.Items != nil {
		s.WriteMore()
		s.WriteObjectField("items")
		a.Items.EncodeSchema(s)
	}

	if a.UniqueItems {
		s.WriteMore()
		s.WriteObjectField("uniqueItems")
		s.WriteBool(a.UniqueItems)
	}

	a.General.EncodeSchema(s)
	s.WriteObjectEnd()
}
