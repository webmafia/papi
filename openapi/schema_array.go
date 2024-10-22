package openapi

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/papi/internal/hasher"
)

var _ Schema = (*Array)(nil)

type Array struct {
	Title       string `tag:"title"`
	Description string `tag:"description"`
	Min         int    `tag:"min"`
	Max         int    `tag:"max"`
	Items       Schema
	Nullable    bool `tag:"flags:nullable"`
	ReadOnly    bool `tag:"flags:readonly"`
	WriteOnly   bool `tag:"flags:writeonly"`
	UniqueItems bool `tag:"flags:unique"`
}

func (sch *Array) GetTitle() string {
	return sch.Title
}

func (sch *Array) encodeSchema(ctx *encoderContext, s *jsoniter.Stream) (err error) {
	if s.Error != nil {
		return s.Error
	}

	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("array")

	if sch.Title != "" {
		s.WriteMore()
		s.WriteObjectField("title")
		s.WriteString(sch.Title)
	}

	if sch.Description != "" {
		s.WriteMore()
		s.WriteObjectField("description")
		s.WriteString(sch.Description)
	}

	if sch.Nullable {
		s.WriteMore()
		s.WriteObjectField("nullable")
		s.WriteBool(sch.Nullable)
	}

	if sch.ReadOnly {
		s.WriteMore()
		s.WriteObjectField("readOnly")
		s.WriteBool(sch.ReadOnly)
	}

	if sch.WriteOnly {
		s.WriteMore()
		s.WriteObjectField("writeOnly")
		s.WriteBool(sch.WriteOnly)
	}

	if sch.Min > 0 {
		s.WriteMore()
		s.WriteObjectField("minItems")
		s.WriteInt(sch.Min)
	}

	if sch.Max > 0 {
		s.WriteMore()
		s.WriteObjectField("maxItems")
		s.WriteInt(sch.Max)
	}

	if sch.Items != nil {
		s.WriteMore()
		s.WriteObjectField("items")
		sch.Items.encodeSchema(ctx, s)
	}

	if sch.UniqueItems {
		s.WriteMore()
		s.WriteObjectField("uniqueItems")
		s.WriteBool(sch.UniqueItems)
	}

	s.WriteObjectEnd()

	if s.Error != nil {
		err = fmt.Errorf("failed to encode array schema: %w", s.Error)
	}

	return
}

func (sch *Array) Hash() uint64 {
	return hasher.Hash(sch)
}
