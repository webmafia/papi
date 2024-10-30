package token

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"sync"
	"unsafe"

	"github.com/webmafia/fast"
	"github.com/webmafia/identifier"
	"github.com/zeebo/blake3"
)

type Generator struct {
	secret []byte
	pool   sync.Pool
}

func NewGenerator(secret []byte) (g *Generator, err error) {
	if len(secret) != 32 {
		return nil, errors.New("token secret must be exactly 32 bytes")
	}

	return &Generator{
		secret: secret,
	}, nil
}

func (g *Generator) CreateToken(payload []byte) (t Token, err error) {
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

	return t, nil
}

func (g *Generator) ValidateToken(t Token) (err error) {
	return g.ValidateTokenBytes(t.bytes())
}

func (g *Generator) ValidateTokenBytes(b []byte) (err error) {
	var signature [32]byte

	if err = g.sign(signature[:0], b[32:]); err != nil {
		return
	}

	if subtle.ConstantTimeCompare(b[:32], signature[:]) == 0 {
		return ErrInvalidAuthToken
	}

	return
}

func (g *Generator) GetValidatedTokenView(b []byte) (tok TokenView, err error) {
	if err = g.ValidateTokenBytes(b); err != nil {
		return
	}

	head := (*sliceHeader)(unsafe.Pointer(&b))

	return tokenView{ptr: unsafe.Add(head.Data, 32)}, nil
}

func (g *Generator) sign(dst, buf []byte) (err error) {
	h := g.acquire()
	defer g.release(h)

	if _, err = h.Write(buf); err != nil {
		return
	}

	h.Sum(dst)
	return
}

func (g *Generator) acquire() *blake3.Hasher {
	if h, ok := g.pool.Get().(*blake3.Hasher); ok {
		return h
	}

	h, _ := blake3.NewKeyed(g.secret)

	return h
}

func (g *Generator) release(h *blake3.Hasher) {
	h.Reset()
	g.pool.Put(h)
}
