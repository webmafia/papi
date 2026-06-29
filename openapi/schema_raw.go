package openapi

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/webmafia/papi/internal/hasher"
)

var _ Schema = (*Raw)(nil)

type Raw struct {
	Title  string `tag:"title"`
	Schema string
}

func (sch *Raw) GetTitle() string {
	return sch.Title
}

func (sch *Raw) encodeSchema(ctx *encoderContext, s *jsoniter.Stream) (err error) {
	if s.Error != nil {
		return s.Error
	}

	s.WriteRaw(sch.Schema)

	if s.Error != nil {
		err = fmt.Errorf("failed to encode raw schema: %w", s.Error)
	}

	return
}

func (sch *Raw) encodeValue(s *jsoniter.Stream, val string) error {
	s.WriteString(val)
	return nil
}

func (sch *Raw) Hash() uint64 {
	return hasher.Hash(sch)
}
