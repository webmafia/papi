package schema

import (
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/fastapi/scan"
)

var _ Schema = (*Object)(nil)

type Object struct {
	General
	Required   []string
	Properties []ObjectProp
}

type ObjectProp struct {
	Name   string
	Schema Schema
}

func (o *Object) ScanTags(tags reflect.StructTag) error {
	return scan.ScanTags(o, tags)
}

// EncodeSchema implements Schema.
func (o Object) EncodeSchema(s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("object")

	if len(o.Required) > 0 {
		s.WriteMore()
		s.WriteObjectField("required")
		s.WriteArrayStart()

		for i := range o.Required {
			if i != 0 {
				s.WriteMore()
			}

			s.WriteString(o.Required[i])
		}

		s.WriteArrayEnd()
	}

	if len(o.Properties) > 0 {
		s.WriteMore()
		s.WriteObjectField("properties")
		s.WriteObjectStart()

		for i, prop := range o.Properties {
			if i != 0 {
				s.WriteMore()
			}

			s.WriteObjectField(prop.Name)
			prop.Schema.EncodeSchema(s)
		}

		s.WriteObjectEnd()
	}

	o.General.EncodeSchema(s)
	s.WriteObjectEnd()
}
