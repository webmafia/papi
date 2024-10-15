package openapi

import (
	"fmt"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var _ Schema = (*Ref)(nil)

type Ref struct {
	Name   string
	Schema Schema
}

func (sch *Ref) GetTitle() string {
	return sch.Name
}

func (sch *Ref) encodeSchema(ctx *encoderContext, s *jsoniter.Stream) (err error) {
	if s.Error != nil {
		return s.Error
	}

	if err = ctx.addRef(sch); err != nil {
		return
	}

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

	if s.Error != nil {
		err = fmt.Errorf("failed to encode reference schema: %w", s.Error)
	}

	return
}

func (sch *Ref) Hash() uint64 {
	return sch.Schema.Hash()
}
