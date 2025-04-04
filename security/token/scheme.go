package token

import (
	"bytes"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/identifier"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/security"
)

var _ security.Scheme = (*Scheme)(nil)

type Scheme struct {
	auth       auth       // Token signing/validation
	tokenStore TokenStore // Token lookup
}

func NewScheme(secret Secret, tokenStore TokenStore) *Scheme {
	return &Scheme{
		auth:       auth{secret: secret},
		tokenStore: tokenStore,
	}
}

// OperationSecurityDocs implements security2.Scheme.
func (s *Scheme) OperationSecurityDocs(permTag string) openapi.SecurityRequirement {
	sec := openapi.SecurityRequirement{
		Name: "token",
	}

	if permTag != "" && permTag != "-" {
		sec.Scopes = []string{permTag}
	}

	return sec
}

// SecurityDocs implements security2.Scheme.
func (s *Scheme) SecurityDocs() openapi.SecurityScheme {
	return openapi.SecurityScheme{
		SchemeName:   "token",
		Type:         "http",
		Description:  "API token",
		Scheme:       "bearer",
		BearerFormat: "base32hex",
	}
}

var tokenPrefix = []byte("Bearer ")

func (s *Scheme) UserRoles(c *fasthttp.RequestCtx) (roles []string, err error) {
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

	return s.tokenStore.UserRoles(c, tok)
}

// Create a token with an optional payload (e.g. a user ID) that will be stored in the token.
// The payload cannot exceed 24 bytes, and will be padded with random bytes.
func (s *Scheme) CreateToken(payload ...[]byte) (t Token, err error) {
	return s.auth.CreateToken(payload...)
}

// Create a token with a specific ID and an optional payload (e.g. a user ID) that will be stored
// in the token. The payload cannot exceed 24 bytes, and will be padded with random bytes.
func (s *Scheme) CreateTokenWithId(id identifier.ID, payload ...[]byte) (t Token, err error) {
	return s.auth.CreateTokenWithId(id, payload...)
}
