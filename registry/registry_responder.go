package registry

import (
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/internal/json"
)

func (r *Registry) Responder(typ reflect.Type) (scan Responder, err error) {

	// Use any existing responder
	if desc, ok := r.describe(typ); ok && desc.Responder != nil {
		return desc.Responder()
	}

	// Use the default responder
	return r.defaultResponder(typ)
}

func (r *Registry) defaultResponder(typ reflect.Type) (Responder, error) {
	enc := json.EncoderOf(typ)

	if typ.Size() == 0 {
		return func(c *fasthttp.RequestCtx, ptr unsafe.Pointer, next func() error) error {
			if err := next(); err != nil {
				return err
			}

			c.SetContentType("application/json")

			return nil
		}, nil
	}

	return func(c *fasthttp.RequestCtx, ptr unsafe.Pointer, next func() error) error {
		if err := next(); err != nil {
			return err
		}

		c.SetContentType("application/json")

		if typ.Size() != 0 {
			s := json.AcquireStream(c.Response.BodyWriter())
			defer json.ReleaseStream(s)

			enc.Encode(ptr, s)
			return s.Flush()
		}

		return nil
	}, nil
}
