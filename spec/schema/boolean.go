package schema

import (
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/fastapi/scan"
)

var _ Schema = (*Boolean)(nil)

type Boolean struct {
	General
}

func (b *Boolean) ScanTags(tags reflect.StructTag) error {
	return scan.ScanTags(b, tags)
}

// EncodeSchema implements Schema.
func (b Boolean) EncodeSchema(s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("boolean")

	b.General.EncodeSchema(s)
	s.WriteObjectEnd()
}
