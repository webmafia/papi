package papi

import (
	"math"
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/webbmaffian/papi/openapi"
	"github.com/webbmaffian/papi/registry"
	"github.com/webmafia/fast"
)

var _ Lister = (*List[struct{}])(nil)

type Lister interface {
	setStream(s *jsoniter.Stream)
	encodeMeta(s *jsoniter.Stream)
}

type List[T any] struct {
	Meta ListMeta

	s   *jsoniter.Stream
	enc jsoniter.ValEncoder
}

type ListMeta struct {
	Total int `json:"total"`
}

func (l *List[T]) setStream(s *jsoniter.Stream) {
	l.s = s
	l.enc = nil
}

func (l *List[T]) encodeMeta(s *jsoniter.Stream) {
	s.WriteObjectStart()
	s.WriteObjectField("total")
	s.WriteInt(l.Meta.Total)
	s.WriteObjectEnd()
}

func (l *List[T]) Write(v *T) {
	// TODO: Don't acquire encoder on-the-fly
	if l.enc == nil {
		l.enc = jsoniter.ConfigFastest.EncoderOf(reflect2.TypeOf(*v))
	} else {
		l.s.WriteMore()
	}

	l.enc.Encode(fast.Noescape(unsafe.Pointer(v)), l.s)
}

func (List[T]) ParamSchema(reg *registry.Registry, tags reflect.StructTag) (schema openapi.Schema, err error) {
	elem, err := reg.Schema(reflect.TypeFor[T]())

	if err != nil {
		return
	}

	sch := &openapi.Object{
		Required: []string{"meta", "items"},
		Properties: []openapi.ObjectProperty{
			{
				Name: "meta",
				Schema: &openapi.Object{
					Title:    "Meta data",
					Required: []string{"total"},
					Properties: []openapi.ObjectProperty{
						{
							Name: "total",
							Schema: &openapi.Integer[int]{
								Max: math.MaxInt,
							},
						},
					},
				},
			},
			{
				Name: "items",
				Schema: &openapi.Array{
					Items: elem,
				},
			},
		},
	}

	if title := elem.GetTitle(); title != "" {
		sch.Title = "List of " + title + " items"
	}

	return sch, nil
}
