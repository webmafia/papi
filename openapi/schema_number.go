package openapi

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/webbmaffian/papi/internal/hasher"
)

var _ Schema = (*Number)(nil)

type Number struct {
	Title       string  `tag:"title"`
	Description string  `tag:"description"`
	Min         float64 `tag:"min"`
	Max         float64 `tag:"max"`
	Nullable    bool    `tag:"flags:nullable"`
	ReadOnly    bool    `tag:"flags:readonly"`
	WriteOnly   bool    `tag:"flags:writeonly"`
}

func (sch *Number) encodeSchema(ctx *encoderContext, s *jsoniter.Stream) (err error) {
	if s.Error != nil {
		return s.Error
	}

	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("number")

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

	if sch.Min >= 0 {
		s.WriteMore()
		s.WriteObjectField("minimum")
		s.WriteFloat64(sch.Min)
	}

	if sch.Max > 0 {
		s.WriteMore()
		s.WriteObjectField("maximum")
		s.WriteFloat64(sch.Max)
	}

	s.WriteObjectEnd()

	if s.Error != nil {
		err = fmt.Errorf("failed to encode number schema: %w", s.Error)
	}

	return
}

func (sch *Number) Hash() uint64 {
	return hasher.Hash(sch)
}
