package openapi

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/papi/internal/constraints"
	"github.com/webmafia/papi/internal/hasher"
)

var _ Schema = (*Number[float32])(nil)

type Number[T constraints.Float] struct {
	Title       string `tag:"title"`
	Description string `tag:"description"`
	Min         T      `tag:"min"`
	Max         T      `tag:"max"`
	Nullable    bool   `tag:"flags:nullable"`
	ReadOnly    bool   `tag:"flags:readonly"`
	WriteOnly   bool   `tag:"flags:writeonly"`
}

func (sch *Number[T]) GetTitle() string {
	return sch.Title
}

func (sch *Number[T]) encodeSchema(ctx *encoderContext, s *jsoniter.Stream) (err error) {
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

	s.WriteMore()
	s.WriteObjectField("minimum")
	s.WriteVal(sch.Min)

	s.WriteMore()
	s.WriteObjectField("maximum")
	s.WriteVal(sch.Max)

	s.WriteObjectEnd()

	if s.Error != nil {
		err = fmt.Errorf("failed to encode number schema: %w", s.Error)
	}

	return
}

func (sch *Number[T]) Hash() uint64 {
	return hasher.Hash(sch)
}
