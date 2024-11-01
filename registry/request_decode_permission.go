package registry

import (
	"bytes"
	"reflect"
	"unsafe"

	"github.com/modern-go/reflect2"
	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/security"
)

func (r *Registry) createPermissionDecoder(typ reflect.Type, perm security.Permission) (scan RequestDecoder, err error) {
	typ2 := reflect2.Type2(typ)
	tokenPrefix := []byte("Bearer ")

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) error {
		rawToken := c.Request.Header.Peek(fasthttp.HeaderAuthorization)
		bearer, ok := bytes.CutPrefix(rawToken, tokenPrefix)

		if !ok {
			return security.ErrInvalidAuthToken
		}

		var tok security.Token

		if err = tok.UnmarshalText(bearer); err != nil {
			return err
		}

		user, err := r.gatekeeper.ValidateToken(c, tok)

		if err != nil {
			return err
		}

		cond, err := r.gatekeeper.GetPolicy(user.UserRoles(), perm)

		if err != nil {
			return err
		}

		if cond != nil {
			typ2.UnsafeSet(p, cond)
		}

		c.SetUserValue("user", user)

		return nil
	}, nil
}
