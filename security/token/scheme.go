package token

import (
	"bytes"
	"iter"
	"reflect"
	"runtime"
	"strings"
	"unsafe"

	"github.com/modern-go/reflect2"
	"github.com/valyala/fasthttp"
	"github.com/webmafia/identifier"
	"github.com/webmafia/papi/internal"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/security"
)

var _ security.Scheme = (*Scheme)(nil)

type Scheme struct {
	auth       auth        // Token signing/validation
	tokenStore TokenStore  // Token lookup
	policies   policyStore // Permissions
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

// OperationSecurityHandler implements security2.Scheme.
func (s *Scheme) OperationSecurityHandler(typ reflect.Type, permTag string, caller *runtime.Func) (handler func(p unsafe.Pointer, c *fasthttp.RequestCtx) error, modTag string, err error) {
	if permTag == "-" {
		return nil, permTag, nil
	}

	perm := Permission(permTag)

	if !perm.HasResource() {
		perm.SetResource(strings.ToLower(internal.CallerTypeFromFunc(caller)))
	}

	if err = s.policies.Register(perm, typ); err != nil {
		return
	}

	typ2 := reflect2.Type2(typ)
	tokenPrefix := []byte("Bearer ")

	return func(p unsafe.Pointer, c *fasthttp.RequestCtx) (err error) {
		rawToken := c.Request.Header.Peek(fasthttp.HeaderAuthorization)
		bearer, ok := bytes.CutPrefix(rawToken, tokenPrefix)

		if !ok {
			return ErrInvalidAuthToken
		}

		var tok Token

		if err = tok.UnmarshalText(bearer); err != nil {
			return
		}

		if err = s.auth.ValidateToken(tok); err != nil {
			return
		}

		user, err := s.tokenStore.Lookup(c, tok)

		if err != nil {
			return
		}

		cond, err := s.GetPolicy(user.UserRoles(), perm)

		if err != nil {
			return err
		}

		if cond != nil {
			typ2.UnsafeSet(p, cond)
		}

		c.SetUserValue("user", user)

		return nil
	}, perm.String(), nil
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

// func (s *Scheme) ValidateToken(t Token) (err error) {
// 	return s.auth.ValidateToken(t)
// }

// Adds a policy. Policies must be added AFTER registering all routes. A policy MIGHT contain either a
// pointer to condition, or a JSON encoded condition as []byte, that will be loaded into a route's policy.
// Any non-matching fields will be ignored. A policy's role + perm combination MUST be unique, or otherwise
// overwritten by the latter. An error will be returned if the permission doesn't exist on any route.
func (s *Scheme) AddPolicy(role string, perm Permission, prio int64, cond ...any) (err error) {
	return s.policies.Add(role, perm, prio, cond...)
}

// Add many policies in bulk. See AddPolicy.
func (s *Scheme) AddPolicies(cb func(add func(role string, perm Permission, prio int64, condJson []byte) error) error) (err error) {
	return s.policies.BatchAdd(cb)
}

// Removes any previously added policy for the role and permission. Does nothing if it never existed.
func (s *Scheme) RemovePolicy(role string, perm Permission) {
	s.policies.Remove(role, perm)
}

// Iterates all added policies.
func (s *Scheme) IteratePolicies() iter.Seq2[PolicyKey, Policy] {
	return s.policies.IteratePolicies()
}

// Iterates all registered permissions. Set inPolicy to iterate permissions either used in policies or not. Default is
// to iterate all regardless it's used in a policy or not.
func (s *Scheme) IteratePermissions(inPolicy ...bool) iter.Seq[Permission] {
	return s.policies.IteratePermissions(inPolicy...)
}

// Any policy matching the route's permission, and one of the user's roles, will be loaded in ascending priority order.
func (s *Scheme) GetPolicy(roles []string, perm Permission) (cond unsafe.Pointer, err error) {
	return s.policies.Get(roles, perm)
}
