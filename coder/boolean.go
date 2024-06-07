package coder

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

var _ ParamCoder[bool] = (*Boolean)(nil)

type Boolean struct {
	General
}

// ScanParam implements ParamCoder.
func (b *Boolean) ScanParam(ptr *bool, str string) error {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "yes", "YES", "Yes":
		*ptr = true
	case "0", "f", "F", "false", "FALSE", "False", "no", "NO", "No":
		*ptr = false
	default:
		return fmt.Errorf("invalid boolean: '%s'", str)
	}

	return nil
}

// EncodeSchema implements Coder.
func (b Boolean) EncodeSchema(s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("boolean")

	b.General.EncodeSchema(s)
	s.WriteObjectEnd()
}
