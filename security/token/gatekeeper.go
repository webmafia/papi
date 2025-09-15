package token

import (
	"bytes"
	"context"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/identifier"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/security"
)

var _ security.Gatekeeper = (*Gatekeeper[struct{}])(nil)
var tokenPrefix = []byte("Bearer ")

type Gatekeeper[T any] struct {
	auth            auth     // Token signing/validation
	store           Store[T] // Token lookup
	optionalPermTag bool
}

func NewGatekeeper[T any](secret Secret, store Store[T], optionalPermTag ...bool) *Gatekeeper[T] {
	g := &Gatekeeper[T]{
		auth:  auth{secret: secret},
		store: store,
	}

	if len(optionalPermTag) > 0 {
		g.optionalPermTag = optionalPermTag[0]
	}

	return g
}

// OptionalPermTag implements security.Gatekeeper.
func (g *Gatekeeper[T]) OptionalPermTag() bool {
	return g.optionalPermTag
}

// OperationSecurityDocs implements security.Gatekeeper.
func (s *Gatekeeper[T]) SecurityRequirement(perm security.Permission) openapi.SecurityRequirement {
	sec := openapi.SecurityRequirement{
		Name: "token",
	}

	if !perm.IsZero() {
		sec.Scopes = []string{perm.String()}
	}

	return sec
}

// SecurityDocs implements security.Gatekeeper.
func (s *Gatekeeper[T]) SecurityScheme() openapi.SecurityScheme {
	return openapi.SecurityScheme{
		SchemeName:   "token",
		Type:         "http",
		Description:  "API token",
		Scheme:       "bearer",
		BearerFormat: "base32hex",
	}
}

// PreRequest implements security.Gatekeeper.
func (g *Gatekeeper[T]) PreRequest(c *fasthttp.RequestCtx) error {
	return nil
}

func (s *Gatekeeper[T]) UserRoles(c *fasthttp.RequestCtx) (roles []string, err error) {
	rawToken := c.Request.Header.Peek(fasthttp.HeaderAuthorization)
	bearer, ok := bytes.CutPrefix(rawToken, tokenPrefix)

	if !ok {
		if cookie := c.Request.Header.Cookie("token"); len(cookie) > 0 {
			bearer = cookie
		} else {
			return nil, security.ErrInvalidAuthToken
		}
	}

	var tok Token

	if err = tok.UnmarshalText(bearer); err != nil {
		return
	}

	if err = s.auth.ValidateToken(tok); err != nil {
		return
	}

	return s.store.UserRoles(c, tok.Id().Int64())
}

func (s *Gatekeeper[T]) CreateAuthCode(ctx context.Context, userId T, expiry time.Duration, cookie bool) (code string, err error) {
	var c OneTimeCode

	if c, err = CreateOneTimeCode(); err != nil {
		return
	}

	code = c.String()
	err = s.store.SaveAuthCode(ctx, userId, code, time.Now().Add(expiry), cookie)
	return
}

// Create a token with an optional payload that will be stored in the token.
// The payload cannot exceed 24 bytes, and will be padded with random bytes.
func (s *Gatekeeper[T]) CreateAccessToken(ctx context.Context, code string, payload ...[]byte) (tok string, cookie bool, err error) {
	userId, cookie, err := s.store.ConsumeAuthCode(ctx, code)

	if err != nil {
		return
	}

	var t Token

	if t, err = s.auth.CreateToken(payload...); err != nil {
		return
	}

	err = s.store.SaveAccessToken(ctx, userId, t.Id().Int64(), cookie)
	return
}

// Consumes an auth code and returns the user ID associated with the auth code
func (s *Gatekeeper[T]) ConsumeAuthCode(ctx context.Context, code string) (userId T, err error) {
	userId, _, err = s.store.ConsumeAuthCode(ctx, code)
	return
}

// Create a token with a specific ID and an optional payload (e.g. a user ID) that will be stored
// in the token. The payload cannot exceed 24 bytes, and will be padded with random bytes.
func (s *Gatekeeper[T]) CreateTokenWithId(id identifier.ID, payload ...[]byte) (string, error) {
	tok, err := s.auth.CreateTokenWithId(id, payload...)

	if err != nil {
		return "", err
	}

	return tok.String(), nil
}
