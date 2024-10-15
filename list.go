package papi

import (
	"math"
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/valyala/fasthttp"
	"github.com/webbmaffian/papi/openapi"
	"github.com/webbmaffian/papi/registry"
	"github.com/webmafia/fast"
)

var _ registry.RequestType = (*List[struct{}])(nil)

type List[T any] struct {
	Meta ListMeta

	s   *jsoniter.Stream
	enc jsoniter.ValEncoder
}

type ListMeta struct {
	Total int `json:"total"`
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

func (List[T]) CreateResponseEncoder(reg *registry.Registry, _ reflect.StructTag, _ []string, handler registry.ResponseEncoder) (registry.ResponseEncoder, error) {
	return func(c *fasthttp.RequestCtx, in, out unsafe.Pointer) error {
		s := reg.JSON().AcquireStream(c.Response.BodyWriter())
		defer reg.JSON().ReleaseStream(s)

		l := (*List[T])(out)
		l.s = s

		s.WriteObjectStart()
		s.WriteObjectField("items")
		s.WriteArrayStart()

		if err := handler(c, in, out); err != nil {
			return err
		}

		s.WriteArrayEnd()
		s.WriteMore()

		s.WriteObjectField("meta")
		s.WriteObjectStart()
		s.WriteObjectField("total")
		s.WriteInt(l.Meta.Total)
		s.WriteObjectEnd()

		s.WriteObjectEnd()

		return s.Flush()
	}, nil
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
