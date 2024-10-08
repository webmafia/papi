package openapi

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var _ Schema = (*Ref)(nil)

type Ref struct {
	Name    string
	Schema  Schema
	written bool
}

func (sch *Ref) encodeSchema(ctx *encoderContext, s *jsoniter.Stream) {
	ctx.addRef(sch)

	s.WriteObjectStart()

	if sch.Name != "" {
		const prefix = "#/components/schemas"
		var b strings.Builder
		b.Grow(len(prefix) + len(sch.Name) + 1)
		b.WriteString(prefix)
		b.WriteByte('/')
		b.WriteString(sch.Name)

		s.WriteObjectField("$ref")
		s.WriteString(b.String())
	}

	s.WriteObjectEnd()
}
