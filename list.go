package fastapi

import jsoniter "github.com/json-iterator/go"

type List[T any] struct {
	Meta ListMeta

	s *jsoniter.Stream
}

type ListMeta struct {
	Total int `json:"total"`
}

func (l *List[T]) setStream(s *jsoniter.Stream) {
	l.s = s
}

func (l List[T]) encodeMeta(s *jsoniter.Stream) {
	s.WriteVal(l.Meta)
}

func (l List[T]) Write(v T) {
	l.s.WriteVal(v)
}

type Lister interface {
	setStream(s *jsoniter.Stream)
	encodeMeta(s *jsoniter.Stream)
}
