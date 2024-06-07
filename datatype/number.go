package datatype

import (
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/fastapi/scan"
)

var _ Type = (*Number)(nil)

type Number struct {
	General
	Min float32 `tag:"min"`
	Max float32 `tag:"max"`
}

func (n *Number) ScanTags(tags reflect.StructTag) error {
	return scan.ScanTags(n, tags)
}

// EncodeSchema implements Type.
func (n Number) EncodeSchema(s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("number")

	if n.Min >= 0 {
		s.WriteMore()
		s.WriteObjectField("minimum")
		s.WriteFloat32Lossy(n.Min)
	}

	if n.Max > 0 {
		s.WriteMore()
		s.WriteObjectField("maximum")
		s.WriteFloat32Lossy(n.Max)
	}

	n.General.EncodeSchema(s)
	s.WriteObjectEnd()
}
