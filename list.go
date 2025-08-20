package papi

import (
	"iter"
	"math"
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
	"github.com/webmafia/papi/internal/json"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
)

var _ registry.TypeDescriber = (*List[struct{}])(nil)

// A streaming list response.
type List[T any] struct {
	meta struct {
		Total int `json:"total"`
	}
	s       *jsoniter.Stream
	enc     jsoniter.ValEncoder
	written bool
}

// Write item to stream.
func (l *List[T]) Write(v *T) error {
	if l.written {
		l.s.WriteMore()
	} else {
		l.written = true
	}

	l.enc.Encode(fast.NoescapeUnsafe(unsafe.Pointer(v)), l.s)
	return l.s.Error
}

// Write item iterator to stream.
func (l *List[T]) WriteAll(it iter.Seq2[*T, error]) (err error) {
	var v *T

	for v, err = range it {
		if err != nil {
			return
		}

		if err = l.Write(v); err != nil {
			return
		}
	}

	return
}

// Set the total number of items that exists (used for e.g. pagination).
func (l *List[T]) SetTotal(i int) {
	l.meta.Total = i
}

// TypeDescription implements registry.TypeDescriber.
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
		Responder: func() (registry.Responder, error) {
			enc := json.EncoderOf(reflect.TypeFor[T]())

			return func(c *fasthttp.RequestCtx, ptr unsafe.Pointer, next func() error) error {
				c.SetContentType("application/json")

				s := json.AcquireStream(c.Response.BodyWriter())
				defer json.ReleaseStream(s)

				l := (*List[T])(ptr)
				l.s = s
				l.enc = enc

				s.WriteObjectStart()
				s.WriteObjectField("items")
				s.WriteArrayStart()

				if err := next(); err != nil {
					return err
				}

				s.WriteArrayEnd()
				s.WriteMore()

				s.WriteObjectField("meta")
				s.WriteObjectStart()
				s.WriteObjectField("total")
				s.WriteInt(l.meta.Total)
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
