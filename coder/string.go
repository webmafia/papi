package coder

import (
	jsoniter "github.com/json-iterator/go"
)

var _ ParamCoder[string] = (*String)(nil)

type String struct {
	General
	Enum    []string `tag:"enum"`
	Format  string   `tag:"format"`
	Pattern string   `tag:"pattern"`
	Min     int      `tag:"min"`
	Max     int      `tag:"max"`
}

// ScanParam implements ParamCoder.
func (s *String) ScanParam(ptr *string, str string) error {
	*ptr = str
	return nil
}

// EncodeSchema implements Coder.
func (str String) EncodeSchema(s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("string")

	if len(str.Enum) > 0 {
		s.WriteMore()
		s.WriteObjectField("enum")
		s.WriteArrayStart()

		for i := range str.Enum {
			if i != 0 {
				s.WriteMore()
			}

			s.WriteString(str.Enum[i])
		}

		s.WriteArrayEnd()
	}

	if str.Format != "" {
		s.WriteMore()
		s.WriteObjectField("format")
		s.WriteString(str.Format)
	}

	if str.Pattern != "" {
		s.WriteMore()
		s.WriteObjectField("pattern")
		s.WriteString(str.Pattern)
	}

	if str.Min > 0 {
		s.WriteMore()
		s.WriteObjectField("minLength")
		s.WriteInt(str.Min)
	}

	if str.Max > 0 {
		s.WriteMore()
		s.WriteObjectField("maxLength")
		s.WriteInt(str.Max)
	}

	str.General.EncodeSchema(s)
	s.WriteObjectEnd()
}
