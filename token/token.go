package token

import (
	"encoding/base32"
	"unsafe"

	"github.com/webmafia/fast"
	"github.com/webmafia/identifier"
)

var encoder = base32.HexEncoding.WithPadding(base32.NoPadding)

var (
	binLen int
	txtLen int
)

func init() {
	binLen = int(unsafe.Sizeof(Token{}))
	txtLen = encoder.EncodedLen(binLen)
}

type Token struct {
	_       [32]byte // Hidden signature
	id      identifier.ID
	payload [24]byte
}

func (t Token) Id() identifier.ID {
	return t.id
}

func (t Token) Payload() [24]byte {
	return t.payload
}

func (t *Token) bytes() []byte {
	return fast.PointerToBytes(t, binLen)
}

func (t Token) String() string {
	buf, _ := t.MarshalText()
	return fast.BytesToString(buf)
}

func (t *Token) FromString(str string) error {
	return t.UnmarshalText(fast.StringToBytes(str))
}

func (t Token) MarshalText() (text []byte, err error) {
	text = make([]byte, txtLen)
	encoder.Encode(text, t.bytes())
	return
}

func (t *Token) UnmarshalText(text []byte) (err error) {
	if len(text) != txtLen {
		return ErrInvalidAuthToken
	}

	_, err = encoder.Decode(t.bytes(), text)
	return
}

func (t Token) MarshalBinary() (data []byte, err error) {
	return t.AppendBinary(make([]byte, 0, binLen))
}

func (t *Token) UnmarshalBinary(data []byte) error {
	if len(data) != binLen {
		return ErrInvalidAuthToken
	}

	copy(t.bytes(), data)
	return nil
}

func (t Token) AppendBinary(b []byte) ([]byte, error) {
	return append(b, t.bytes()...), nil
}

func (t Token) AppendText(b []byte) ([]byte, error) {
	return encoder.AppendEncode(b, t.bytes()), nil
}
