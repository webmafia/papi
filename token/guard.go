package token

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"sync"

	"github.com/webmafia/fast"
	"github.com/webmafia/identifier"
	"github.com/zeebo/blake3"
)

type Gatekeeper struct {
	secret Secret
	pool   sync.Pool
	store  TokenStore
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

func (g *Gatekeeper) CreateToken(ctx context.Context, payload []byte) (t Token, err error) {
	t = Token{
		id: identifier.Generate(),
	}

	if _, err = rand.Read(fast.NoescapeBytes(t.payload[:])); err != nil {
		return
	}

	copy(t.payload[:], payload)
	b := t.bytes()

	if err = g.sign(b[:0], b[32:]); err != nil {
		return
	}

	if err = g.store.Store(ctx, t); err != nil {
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
	h := g.acquire()
	defer g.release(h)

	if _, err = h.Write(buf); err != nil {
		return
	}

	h.Sum(dst)
	return
}

func (g *Gatekeeper) acquire() *blake3.Hasher {
	if h, ok := g.pool.Get().(*blake3.Hasher); ok {
		return h
	}

	h, _ := blake3.NewKeyed(g.secret[:])

	return h
}

func (g *Gatekeeper) release(h *blake3.Hasher) {
	h.Reset()
	g.pool.Put(h)
}
