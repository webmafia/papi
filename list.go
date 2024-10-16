package papi

import (
	"math"
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"github.com/webbmaffian/papi/openapi"
	"github.com/webbmaffian/papi/registry"
	"github.com/webmafia/fast"
)

var _ registry.TypeDescriber = (*List[struct{}])(nil)

type List[T any] struct {
	Meta ListMeta

	s       *jsoniter.Stream
	enc     jsoniter.ValEncoder
	written bool
}

type ListMeta struct {
	Total int `json:"total"`
}

func (l *List[T]) Write(v *T) {
	if l.written {
		l.s.WriteMore()
	} else {
		l.written = true
	}

	l.enc.Encode(fast.Noescape(unsafe.Pointer(v)), l.s)
}

func (List[T]) TypeDescription(reg *registry.Registry) registry.TypeDescription {
	return registry.TypeDescription{
		Schema: func(_ reflect.StructTag) (schema openapi.Schema, err error) {
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
		},
		Handler: func(_ reflect.StructTag, handler registry.Handler) (registry.Handler, error) {
			enc := reg.JSON().EncoderOf(reflect.TypeFor[T]())

			return func(c *fasthttp.RequestCtx, in, out unsafe.Pointer) error {
				c.SetContentType("application/json")

				s := reg.JSON().AcquireStream(c.Response.BodyWriter())
				defer reg.JSON().ReleaseStream(s)

				l := (*List[T])(out)
				l.s = s
				l.enc = enc

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

				l.s = nil
				l.enc = nil
				l.written = false

				return s.Flush()
			}, nil
		},
	}
}
