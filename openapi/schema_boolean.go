package openapi

import (
	"fmt"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/papi/internal/hasher"
)

var _ Schema = (*Boolean)(nil)

type Boolean struct {
	Title       string `tag:"title"`
	Description string `tag:"description"`
	Default     string `tag:"default"`
	Nullable    bool   `tag:"flags:nullable"`
	ReadOnly    bool   `tag:"flags:readonly"`
	WriteOnly   bool   `tag:"flags:writeonly"`
}

func (sch *Boolean) GetTitle() string {
	return sch.Title
}

func (sch *Boolean) encodeSchema(ctx *encoderContext, s *jsoniter.Stream) (err error) {
	if s.Error != nil {
		return s.Error
	}

	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("boolean")

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

	if sch.Default != "" {
		s.WriteMore()
		s.WriteObjectField("default")
		sch.encodeValue(s, sch.Default)
	}

	s.WriteObjectEnd()

	if s.Error != nil {
		err = fmt.Errorf("failed to encode boolean schema: %w", s.Error)
	}

	return
}

func (sch *Boolean) encodeValue(s *jsoniter.Stream, val string) error {
	val = strings.ToLower(val)

	if val == "true" {
		s.WriteTrue()
	} else if val == "false" {
		s.WriteFalse()
	} else {
		return fmt.Errorf("invalid boolean: '%s'", val)
	}

	return nil
}

func (sch *Boolean) Hash() uint64 {
	return hasher.Hash(sch)
}
