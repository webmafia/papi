package coder

import (
	"reflect"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/fastapi/scan"
)

var (
	_ ParamCoder[int]    = (*IntegerInt)(nil)
	_ ParamCoder[int8]   = (*IntegerInt8)(nil)
	_ ParamCoder[int16]  = (*IntegerInt16)(nil)
	_ ParamCoder[int32]  = (*IntegerInt32)(nil)
	_ ParamCoder[int64]  = (*IntegerInt64)(nil)
	_ ParamCoder[uint]   = (*IntegerUint)(nil)
	_ ParamCoder[uint8]  = (*IntegerUint8)(nil)
	_ ParamCoder[uint16] = (*IntegerUint16)(nil)
	_ ParamCoder[uint32] = (*IntegerUint32)(nil)
	_ ParamCoder[uint64] = (*IntegerUint64)(nil)
)

type Integer struct {
	General
	Min int `tag:"min"`
	Max int `tag:"max"`
}

func (i *Integer) ScanTags(tags reflect.StructTag) error {
	return scan.ScanTags(i, tags)
}

// EncodeSchema implements Coder.
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

type IntegerInt Integer

// ScanParam implements ParamCoder.
func (i *IntegerInt) ScanParam(ptr *int, str string) (err error) {
	v, err := strconv.ParseInt(str, 10, 0)
	if err == nil {
		*ptr = int(v)
	}
	return
}

type IntegerInt8 Integer

// ScanParam implements ParamCoder.
func (i *IntegerInt8) ScanParam(ptr *int8, str string) (err error) {
	v, err := strconv.ParseInt(str, 10, 8)
	if err == nil {
		*ptr = int8(v)
	}
	return
}

type IntegerInt16 Integer

// ScanParam implements ParamCoder.
func (i *IntegerInt16) ScanParam(ptr *int16, str string) (err error) {
	v, err := strconv.ParseInt(str, 10, 16)
	if err == nil {
		*ptr = int16(v)
	}
	return
}

type IntegerInt32 Integer

// ScanParam implements ParamCoder.
func (i *IntegerInt32) ScanParam(ptr *int32, str string) (err error) {
	v, err := strconv.ParseInt(str, 10, 32)
	if err == nil {
		*ptr = int32(v)
	}
	return
}

type IntegerInt64 Integer

// ScanParam implements ParamCoder.
func (i *IntegerInt64) ScanParam(ptr *int64, str string) (err error) {
	v, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		*ptr = v
	}
	return
}

type IntegerUint Integer

// ScanParam implements ParamCoder.
func (i *IntegerUint) ScanParam(ptr *uint, str string) (err error) {
	v, err := strconv.ParseUint(str, 10, 0)
	if err == nil {
		*ptr = uint(v)
	}
	return
}

type IntegerUint8 Integer

// ScanParam implements ParamCoder.
func (i *IntegerUint8) ScanParam(ptr *uint8, str string) (err error) {
	v, err := strconv.ParseUint(str, 10, 8)
	if err == nil {
		*ptr = uint8(v)
	}
	return
}

type IntegerUint16 Integer

// ScanParam implements ParamCoder.
func (i *IntegerUint16) ScanParam(ptr *uint16, str string) (err error) {
	v, err := strconv.ParseUint(str, 10, 16)
	if err == nil {
		*ptr = uint16(v)
	}
	return
}

type IntegerUint32 Integer

// ScanParam implements ParamCoder.
func (i *IntegerUint32) ScanParam(ptr *uint32, str string) (err error) {
	v, err := strconv.ParseUint(str, 10, 32)
	if err == nil {
		*ptr = uint32(v)
	}
	return
}

type IntegerUint64 Integer

// ScanParam implements ParamCoder.
func (i *IntegerUint64) ScanParam(ptr *uint64, str string) (err error) {
	v, err := strconv.ParseUint(str, 10, 64)
	if err == nil {
		*ptr = v
	}
	return
}
