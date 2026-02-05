package openapi

import (
	"errors"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/papi/internal/hasher"
)

var _ Schema = (*Object)(nil)

type Object struct {
	Title                string `tag:"title"`
	Description          string `tag:"description"`
	Required             []string
	Properties           []ObjectProperty
	AdditionalProperties Schema
	Nullable             bool `tag:"flags:nullable"`
	ReadOnly             bool `tag:"flags:readonly"`
	WriteOnly            bool `tag:"flags:writeonly"`
}

type ObjectProperty struct {
	Name   string
	Schema Schema
}

func (sch *Object) GetTitle() string {
	return sch.Title
}

func (sch *Object) encodeSchema(ctx *encoderContext, s *jsoniter.Stream) (err error) {
	if s.Error != nil {
		return s.Error
	}

	s.WriteObjectStart()

	s.WriteObjectField("type")
	s.WriteString("object")

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

	if len(sch.Required) > 0 {
		s.WriteMore()
		s.WriteObjectField("required")
		s.WriteArrayStart()

		for i := range sch.Required {
			if i != 0 {
				s.WriteMore()
			}

			s.WriteString(sch.Required[i])
		}

		s.WriteArrayEnd()
	}

	if len(sch.Properties) > 0 {
		s.WriteMore()
		s.WriteObjectField("properties")
		s.WriteObjectStart()

		for i := range sch.Properties {
			if i != 0 {
				s.WriteMore()
			}

			s.WriteObjectField(sch.Properties[i].Name)

			if err = encodeSchema(ctx, s, sch.Properties[i].Schema); err != nil {
				return
			}
		}

		s.WriteObjectEnd()
	}

	if sch.AdditionalProperties != nil {
		s.WriteMore()
		s.WriteObjectField("additionalProperties")

		if err = sch.AdditionalProperties.encodeSchema(ctx, s); err != nil {
			return
		}
	}

	s.WriteObjectEnd()

	if s.Error != nil {
		err = fmt.Errorf("failed to encode object schema: %w", s.Error)
	}

	return
}

func (sch *Object) encodeValue(s *jsoniter.Stream, val string) error {
	return errors.New("default objects not supported")
}

func (sch *Object) Hash() uint64 {
	return hasher.Hash(sch)
}
