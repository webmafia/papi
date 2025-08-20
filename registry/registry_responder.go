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

	return func(c *fasthttp.RequestCtx, ptr unsafe.Pointer, next func() error) error {
		if err := next(); err != nil {
			return err
		}

		c.SetContentType("application/json")

		s := json.AcquireStream(c.Response.BodyWriter())
		defer json.ReleaseStream(s)

		enc.Encode(ptr, s)
		return s.Flush()
	}, nil
}
