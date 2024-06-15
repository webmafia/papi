package coder

import (
	"fmt"
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/fastapi/scan"
)

var _ ParamCoder[bool] = Boolean{}

type Boolean struct{}

// ScanParam implements ParamCoder.
func (Boolean) ScanParam(ptr *bool, str string) error {
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
func (Boolean) Encoder(tag reflect.StructTag) func(s *jsoniter.Stream) {
	var tags General

	scan.ScanTags(&tags, tag)

	return func(s *jsoniter.Stream) {
		s.WriteObjectStart()

		s.WriteObjectField("type")
		s.WriteString("boolean")

		tags.EncodeSchema(s)
		s.WriteObjectEnd()
	}
}
