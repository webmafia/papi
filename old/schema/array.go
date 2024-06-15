package schema

import (
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/fastapi/scan"
)

var _ Schema = (*Array)(nil)

type Array struct {
	General
	Min         int `tag:"min"`
	Max         int `tag:"max"`
	Items       Schema
	UniqueItems bool `tag:"unique"`
}

func (a *Array) ScanTags(tags reflect.StructTag) error {
	return scan.ScanTags(a, tags)
}

// EncodeSchema implements Schema.
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
