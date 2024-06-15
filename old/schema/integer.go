package schema

import (
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/fastapi/scan"
)

var _ Schema = (*Integer)(nil)

type Integer struct {
	General
	Min int `tag:"min"`
	Max int `tag:"max"`
}

func (i *Integer) ScanTags(tags reflect.StructTag) error {
	return scan.ScanTags(i, tags)
}

// EncodeSchema implements Schema.
func (i Integer) EncodeSchema(s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("integer")

	if i.Min >= 0 {
		s.WriteMore()
		s.WriteObjectField("minimum")
		s.WriteInt(i.Min)
	}

	if i.Max > 0 {
		s.WriteMore()
		s.WriteObjectField("maximum")
		s.WriteInt(i.Max)
	}

	i.General.EncodeSchema(s)
	s.WriteObjectEnd()
}
