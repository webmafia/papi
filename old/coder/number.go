package coder

import (
	"reflect"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/fastapi/scan"
)

var (
	_ ParamCoder[float32] = (*NumberFloat32)(nil)
	_ ParamCoder[float64] = (*NumberFloat64)(nil)
)

type Number struct {
	General
	Min float64 `tag:"min"`
	Max float64 `tag:"max"`
}

func (n *Number) ScanTags(tags reflect.StructTag) error {
	return scan.ScanTags(n, tags)
}

// EncodeSchema implements Coder.
func (n Number) EncodeSchema(s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("number")

	if n.Min >= 0 {
		s.WriteMore()
		s.WriteObjectField("minimum")
		s.WriteFloat64Lossy(n.Min)
	}

	if n.Max > 0 {
		s.WriteMore()
		s.WriteObjectField("maximum")
		s.WriteFloat64Lossy(n.Max)
	}

	n.General.EncodeSchema(s)
	s.WriteObjectEnd()
}

type NumberFloat32 Number

// ScanParam implements ParamCoder.
func (i *NumberFloat32) ScanParam(ptr *float32, str string) (err error) {
	v, err := strconv.ParseFloat(str, 32)
	if err == nil {
		*ptr = float32(v)
	}
	return
}

type NumberFloat64 Number

// ScanParam implements ParamCoder.
func (i *NumberFloat64) ScanParam(ptr *float64, str string) (err error) {
	v, err := strconv.ParseFloat(str, 64)
	if err == nil {
		*ptr = v
	}
	return
}
