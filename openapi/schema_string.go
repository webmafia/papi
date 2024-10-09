package openapi

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/webbmaffian/papi/internal/hasher"
)

var _ Schema = (*String)(nil)

type String struct {
	Title       string   `tag:"title"`
	Description string   `tag:"description"`
	Enum        []string `tag:"enum"`
	Format      string   `tag:"format"`
	Pattern     string   `tag:"pattern"`
	Min         int      `tag:"min"`
	Max         int      `tag:"max"`
	Nullable    bool     `tag:"flags:nullable"`
	ReadOnly    bool     `tag:"flags:readonly"`
	WriteOnly   bool     `tag:"flags:writeonly"`
}

func (sch *String) encodeSchema(ctx *encoderContext, s *jsoniter.Stream) (err error) {
	if s.Error != nil {
		return s.Error
	}

	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("string")

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

	if len(sch.Enum) > 0 {
		s.WriteMore()
		s.WriteObjectField("enum")
		s.WriteArrayStart()

		for i := range sch.Enum {
			if i != 0 {
				s.WriteMore()
			}

			s.WriteString(sch.Enum[i])
		}

		s.WriteArrayEnd()
	}

	if sch.Format != "" {
		s.WriteMore()
		s.WriteObjectField("format")
		s.WriteString(sch.Format)
	}

	if sch.Pattern != "" {
		s.WriteMore()
		s.WriteObjectField("pattern")
		s.WriteString(sch.Pattern)
	}

	if sch.Min > 0 {
		s.WriteMore()
		s.WriteObjectField("minLength")
		s.WriteInt(sch.Min)
	}

	if sch.Max > 0 {
		s.WriteMore()
		s.WriteObjectField("maxLength")
		s.WriteInt(sch.Max)
	}

	s.WriteObjectEnd()

	if s.Error != nil {
		err = fmt.Errorf("failed to encode string schema: %w", s.Error)
	}

	return
}

func (sch *String) Hash() uint64 {
	return hasher.Hash(sch)
}
