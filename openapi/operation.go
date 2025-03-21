package openapi

import (
	jsoniter "github.com/json-iterator/go"
)

type Operation struct {
	Id          string
	Method      string
	Summary     string
	Description string
	Security    []SecurityRequirement
	Parameters  []Parameter
	RequestBody Schema
	Response    Schema
	Tags        []Tag
}

func (op *Operation) JsonEncode(ctx *encoderContext, s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("summary")
	s.WriteString(op.Summary)

	if op.Description != "" {
		s.WriteMore()
		s.WriteObjectField("description")
		s.WriteString(op.Description)
	}

	s.WriteMore()
	s.WriteObjectField("operationId")
	s.WriteString(op.Id)

	if len(op.Security) > 0 {
		s.WriteMore()
		s.WriteObjectField("security")
		s.WriteArrayStart()

		for i := range op.Security {
			if i != 0 {
				s.WriteMore()
			}

			s.WriteObjectStart()

			s.WriteObjectField(op.Security[i].Name)
			op.Security[i].JsonEncode(s)

			s.WriteObjectEnd()
		}

		s.WriteArrayEnd()
	}

	if len(op.Parameters) > 0 {
		s.WriteMore()
		s.WriteObjectField("parameters")

		s.WriteArrayStart()

		for i := range op.Parameters {
			if i != 0 {
				s.WriteMore()
			}

			op.Parameters[i].JsonEncode(ctx, s)
		}

		s.WriteArrayEnd()
	}

	if len(op.Tags) > 0 {
		s.WriteMore()
		s.WriteObjectField("tags")
		s.WriteArrayStart()

		for i := range op.Tags {
			if i != 0 {
				s.WriteMore()
			}

			s.WriteString(op.Tags[i].Name)
			ctx.addTag(op.Tags[i])
		}

		s.WriteArrayEnd()
	}

	if op.Method != "get" && op.Method != "GET" {
		s.WriteMore()
		s.WriteObjectField("requestBody")
		s.WriteObjectStart()

		s.WriteObjectField("content")
		s.WriteObjectStart()

		s.WriteObjectField("application/json")
		s.WriteObjectStart()

		s.WriteObjectField("schema")
		encodeSchema(ctx, s, op.RequestBody)

		s.WriteObjectEnd()
		s.WriteObjectEnd()
		s.WriteObjectEnd()
	}

	s.WriteMore()
	s.WriteObjectField("responses")
	s.WriteObjectStart()

	s.WriteObjectField("200")
	s.WriteObjectStart()

	s.WriteObjectField("description")

	if op.Response != nil {
		if title := op.Response.GetTitle(); title != "" {
			s.WriteString(title)
		} else {
			s.WriteString("Response")
		}
	} else {
		s.WriteString("Response")
	}

	s.WriteMore()
	s.WriteObjectField("content")
	s.WriteObjectStart()

	if custom, ok := op.Response.(*Custom); ok {
		s.WriteObjectField(custom.ContentType)
	} else {
		s.WriteObjectField("application/json")
	}

	s.WriteObjectStart()

	s.WriteObjectField("schema")
	encodeSchema(ctx, s, op.Response)

	s.WriteObjectEnd()

	s.WriteObjectEnd()

	s.WriteObjectEnd()

	s.WriteObjectEnd()

	s.WriteObjectEnd()
}

func encodeSchema(ctx *encoderContext, s *jsoniter.Stream, sch Schema) {
	if sch != nil {
		sch.encodeSchema(ctx, s)
	} else {
		s.WriteEmptyObject()
	}
}
