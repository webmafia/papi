package security

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"iter"
	"reflect"
	"sync"
	"unsafe"

	"github.com/webmafia/fast"
	"github.com/webmafia/identifier"
	"github.com/zeebo/blake3"
)

type Gatekeeper struct {
	secret   Secret
	pool     sync.Pool
	store    TokenStore
	policies policyStore
}

func NewGatekeeper(secret Secret, store TokenStore) (g *Gatekeeper, err error) {
	if len(secret) != 32 {
		return nil, errors.New("token secret must be exactly 32 bytes")
	}

	return &Gatekeeper{
		secret: secret,
		store:  store,
	}, nil
}

// Create a token with an optional payload (e.g. a user ID) that will be stored in the token.
// The payload cannot exceed 24 bytes, and will be padded with random bytes.
func (g *Gatekeeper) CreateToken(ctx context.Context, payload ...[]byte) (t Token, err error) {
	t = Token{
		id: identifier.Generate(),
	}

	var payloadSize int

	if len(payload) > 0 {
		if len(payload[0]) > 24 {
			err = errors.New("payload cannot exceed 24 bytes")
			return
		}

		payloadSize = copy(t.payload[:], payload[0])
	}

	if _, err = rand.Read(fast.NoescapeBytes(t.payload[payloadSize:])); err != nil {
		return
	}

	b := t.bytes()

	if err = g.sign(b[:0], b[32:]); err != nil {
		return
	}

	if err = g.store.Insert(ctx, t); err != nil {
		return
	}

	return t, nil
}

func (g *Gatekeeper) ValidateToken(ctx context.Context, t Token) (user User, err error) {
	if err = g.validateTokenBytes(fast.NoescapeBytes(t.bytes())); err != nil {
		return
	}

	return g.store.Lookup(ctx, t)
}

func (g *Gatekeeper) validateTokenBytes(b []byte) (err error) {
	var signature [32]byte

	if err = g.sign(signature[:0], b[32:]); err != nil {
		return
	}

	if subtle.ConstantTimeCompare(b[:32], signature[:]) == 0 {
		return ErrInvalidAuthToken
	}

	return
}

func (g *Gatekeeper) sign(dst, buf []byte) (err error) {
	h := g.acquireSigner()
	defer g.releaseSigner(h)

	if _, err = h.Write(buf); err != nil {
		return
	}

	h.Sum(dst)
	return
}

func (g *Gatekeeper) acquireSigner() *blake3.Hasher {
	if h, ok := g.pool.Get().(*blake3.Hasher); ok {
		return h
	}

	// We can safely ignore this error, as we can guarantee that the key is the right size
	h, _ := blake3.NewKeyed(g.secret[:])

	return h
}

func (g *Gatekeeper) releaseSigner(h *blake3.Hasher) {
	h.Reset()
	g.pool.Put(h)
}

// Registers a permission. This should NOT be called manually, as it's called automatically when
// registering routes.
func (g *Gatekeeper) RegisterPermission(perm Permission, typ reflect.Type) (err error) {
	return g.policies.Register(perm, typ)
}

// Adds a policy. Policies must be added AFTER registering all routes. A policy MIGHT contain a JSON
// encoded condition, that will be loaded into a route's policy. Any non-matching fields will be ignored.
// A policy's role + perm combination MUST be unique, or otherwise overwritten by the latter. An error
// will be returned if the permission doesn't exist on any route.
func (g *Gatekeeper) AddPolicy(role string, perm Permission, prio int64, condJson []byte) (err error) {
	return g.policies.Add(role, perm, prio, condJson)
}

// Add many policies in bulk. See AddPolicy.
func (g *Gatekeeper) AddPolicies(cb func(add func(role string, perm Permission, prio int64, condJson []byte) error) error) (err error) {
	return g.policies.BatchAdd(cb)
}

// Removes any previously added policy for the role and permission. Does nothing if it never existed.
func (g *Gatekeeper) RemovePolicy(role string, perm Permission) {
	g.policies.Remove(role, perm)
}

// Iterates all added policies.
func (g *Gatekeeper) IteratePolicies() iter.Seq2[PolicyKey, Policy] {
	return g.policies.IteratePolicies()
}

// Iterates all registered permissions. Set inPolicy to iterate permissions either used in policies or not. Default is
// to iterate all regardless it's used in a policy or not.
func (g *Gatekeeper) IteratePermissions(inPolicy ...bool) iter.Seq[Permission] {
	return g.policies.IteratePermissions(inPolicy...)
}

// Any policy matching the route's permission, and one of the user's roles, will be loaded in ascending priority order.
func (g *Gatekeeper) GetPolicy(roles []string, perm Permission) (cond unsafe.Pointer, err error) {
	return g.policies.Get(roles, perm)
}
