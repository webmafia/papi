package papi

import (
	"encoding/json"
	"errors"
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
)

var (
	_ registry.TypeDescriber = (*RawJSON)(nil)
	_ json.Marshaler         = (RawJSON)(nil)
	_ json.Unmarshaler       = (*RawJSON)(nil)
)

type RawJSON []byte

// TypeDescription implements registry.TypeDescriber.
func (RawJSON) TypeDescription(reg *registry.Registry) registry.TypeDescription {
	return registry.TypeDescription{
		Schema: func(tags reflect.StructTag) (schema openapi.Schema, err error) {
			return nil, nil
		},
		Handler: func(tags reflect.StructTag, handler registry.Handler) (registry.Handler, error) {
			return func(c *fasthttp.RequestCtx, in, out unsafe.Pointer) (err error) {
				if err = handler(c, in, out); err != nil {
					return
				}

				o := (*RawJSON)(out)

				if len(*o) > 0 {
					c.Response.SetBody(*o)
				} else {
					c.Response.SetBodyString("{}")
				}

				return
			}, nil
		},
		Decoder: func(tags reflect.StructTag) (registry.Decoder, error) {
			return func(p unsafe.Pointer, s string) error {
				b := fast.StringToBytes(s)

				if !jsoniter.Valid(b) {
					return errors.New("invalid JSON")
				}

				i := (*RawJSON)(p)
				*i = b
				return nil
			}, nil
		},
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *RawJSON) UnmarshalJSON(b []byte) error {
	*r = append((*r)[:0], b...)
	return nil
}

// MarshalJSON implements json.Marshaler.
func (r RawJSON) MarshalJSON() ([]byte, error) {
	return r, nil
}
