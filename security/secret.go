package security

import (
	"crypto/rand"
	"fmt"

	"github.com/webmafia/fast"
)

const secretLen = 32

var secretEncLen = encoder.EncodedLen(secretLen)

func SecretFromString(str string) (s Secret, err error) {
	err = s.FromString(str)
	return
}

func GenerateSecret() (s Secret, err error) {
	_, err = rand.Read(s[:])
	return
}

type Secret [secretLen]byte

func (s Secret) String() string {
	buf, _ := s.MarshalText()
	return fast.BytesToString(buf)
}

func (s *Secret) FromString(str string) error {
	return s.UnmarshalText(fast.StringToBytes(str))
}

func (s Secret) MarshalText() (text []byte, err error) {
	text = make([]byte, secretEncLen)
	encoder.Encode(text, s[:])
	return
}

func (s *Secret) UnmarshalText(text []byte) (err error) {
	if len(text) != secretEncLen {
		return fmt.Errorf("token secret must be exactly %d bytes, endoded to a %d characters string", secretLen, secretEncLen)
	}

	_, err = encoder.Decode((*s)[:], text)
	return
}

func (t Secret) MarshalBinary() (data []byte, err error) {
	return t.AppendBinary(make([]byte, 0, binLen))
}

func (s *Secret) UnmarshalBinary(data []byte) error {
	if len(data) != secretLen {
		return fmt.Errorf("token secret must be exactly %d bytes", secretLen)
	}

	copy((*s)[:], data)
	return nil
}

func (s Secret) AppendBinary(b []byte) ([]byte, error) {
	return append(b, s[:]...), nil
}

func (s Secret) AppendText(b []byte) ([]byte, error) {
	return encoder.AppendEncode(b, s[:]), nil
}
