package datatype

import jsoniter "github.com/json-iterator/go"

type General struct {
	name      string
	Nullable  bool `tag:"nullable"`
	ReadOnly  bool `tag:"readOnly"`
	WriteOnly bool `tag:"writeOnly"`
}

func (g General) Name() string {
	return g.name
}

func (g *General) SetName(name string) {
	g.name = name
}

func (g General) EncodeSchema(s *jsoniter.Stream) {
	if g.Nullable {
		s.WriteMore()
		s.WriteObjectField("nullable")
		s.WriteTrue()
	}

	if g.ReadOnly {
		s.WriteMore()
		s.WriteObjectField("readOnly")
		s.WriteTrue()
	}

	if g.WriteOnly {
		s.WriteMore()
		s.WriteObjectField("writeOnly")
		s.WriteTrue()
	}
}
