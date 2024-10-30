package openapi

import (
	jsoniter "github.com/json-iterator/go"
)

type Operation struct {
	Id          string
	Method      string
	Summary     string
	Description string
	Security    Security
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

	s.WriteMore()
	op.Security.JsonEncode(ctx, s)

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

	if op.RequestBody != nil {
		s.WriteMore()
		s.WriteObjectField("requestBody")
		s.WriteObjectStart()

		s.WriteObjectField("content")
		s.WriteObjectStart()

		s.WriteObjectField("application/json")
		s.WriteObjectStart()

		s.WriteObjectField("schema")
		op.RequestBody.encodeSchema(ctx, s)

		s.WriteObjectEnd()

		s.WriteObjectEnd()

		s.WriteObjectEnd()
	}

	if op.Response != nil {
		s.WriteMore()
		s.WriteObjectField("responses")
		s.WriteObjectStart()

		s.WriteObjectField("200")
		s.WriteObjectStart()

		s.WriteObjectField("description")

		if title := op.Response.GetTitle(); title != "" {
			s.WriteString(title)
		} else {
			s.WriteString("Response")
		}

		s.WriteMore()
		s.WriteObjectField("content")
		s.WriteObjectStart()

		s.WriteObjectField("application/json")
		s.WriteObjectStart()

		s.WriteObjectField("schema")
		op.Response.encodeSchema(ctx, s)

		s.WriteObjectEnd()

		s.WriteObjectEnd()

		s.WriteObjectEnd()

		s.WriteObjectEnd()
	}

	s.WriteObjectEnd()
}
