package token

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"sync"

	"github.com/webmafia/fast"
	"github.com/webmafia/hexid"
	"github.com/webmafia/papi/security"
	"github.com/zeebo/blake3"
)

type auth struct {
	secret Secret
	pool   sync.Pool
}

// Create a token with an optional payload (e.g. a user ID) that will be stored in the token.
// The payload cannot exceed 24 bytes, and will be padded with random bytes.
func (g *auth) CreateToken(payload ...[]byte) (t Token, err error) {
	return g.CreateTokenWithId(hexid.Generate(), payload...)
}

// Create a token with a specific ID and an optional payload (e.g. a user ID) that will be stored
// in the token. The payload cannot exceed 24 bytes, and will be padded with random bytes.
func (g *auth) CreateTokenWithId(id hexid.ID, payload ...[]byte) (t Token, err error) {
	t = Token{
		id: id,
	}

	var payloadSize int

	if len(payload) > 0 {
		if len(payload[0]) > 24 {
			err = errors.New("payload cannot exceed 24 bytes")
			return
		}

		payloadSize = copy(t.payload[:], payload[0])
	}

	if _, err = rand.Read(fast.Noescape(t.payload[payloadSize:])); err != nil {
		return
	}

	b := t.bytes()

	if err = g.sign(b[:0], b[32:]); err != nil {
		return
	}

	return t, nil
}

func (g *auth) ValidateToken(t Token) (err error) {
	return g.validateTokenBytes(fast.Noescape(t.bytes()))
}

func (g *auth) validateTokenBytes(b []byte) (err error) {
	var signature [32]byte

	if err = g.sign(signature[:0], b[32:]); err != nil {
		return
	}

	if subtle.ConstantTimeCompare(b[:32], signature[:]) == 0 {
		return security.ErrInvalidAuthToken
	}

	return
}

func (g *auth) sign(dst, buf []byte) (err error) {
	h := g.acquireSigner()
	defer g.releaseSigner(h)

	if _, err = h.Write(buf); err != nil {
		return
	}

	h.Sum(dst)
	return
}

func (g *auth) acquireSigner() *blake3.Hasher {
	if h, ok := g.pool.Get().(*blake3.Hasher); ok {
		return h
	}

	// We can safely ignore this error, as we can guarantee that the key is the right size
	h, _ := blake3.NewKeyed(g.secret[:])

	return h
}

func (g *auth) releaseSigner(h *blake3.Hasher) {
	h.Reset()
	g.pool.Put(h)
}
