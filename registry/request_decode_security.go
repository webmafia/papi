package registry

import (
	"bytes"
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/token"
)

func (r *Registry) createSecurityDecoder(typ reflect.Type, action, resource string) (scan RequestDecoder, err error) {
	tokenPrefix := []byte("Bearer ")

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		rawToken := c.Request.Header.Peek(fasthttp.HeaderAuthorization)
		bearer, ok := bytes.CutPrefix(rawToken, tokenPrefix)

		if !ok {
			if action == "" {
				return nil
			}

			return token.ErrInvalidAuthToken
		}

		tok, err := r.tokGen.GetValidatedTokenView(bearer)

		if err != nil {
			return err
		}

		c.SetUserValue("token", tok)

		if action != "" {
			// TODO: Get role from token, then check policy.
		}

		return nil
	}, nil
}
