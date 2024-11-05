package token

import (
	"crypto/rand"
	"fmt"

	"github.com/webmafia/fast"
)

const codeTxtLen = 32
const codeFmtLen = codeTxtLen + codeTxtLen/4 - 1

func CreateOneTimeCode() (otc OneTimeCode, err error) {
	_, err = rand.Read(otc[:])
	return
}

type OneTimeCode [20]byte

func (t OneTimeCode) String() string {
	buf, _ := t.MarshalText()
	return fast.BytesToString(buf)
}

func (t *OneTimeCode) FromString(str string) error {
	return t.UnmarshalText(fast.StringToBytes(str))
}

func (t OneTimeCode) MarshalText() (text []byte, err error) {
	text = make([]byte, codeTxtLen, codeFmtLen)
	encoder.Encode(text, t[:])
	text = fmt.Appendf(text[:0],
		"%s-%s-%s-%s-%s-%s-%s-%s",
		text[:4],
		text[4:8],
		text[8:12],
		text[12:16],
		text[16:20],
		text[20:24],
		text[24:28],
		text[28:],
	)
	return
}

func (t *OneTimeCode) UnmarshalText(text []byte) (err error) {
	if len(text) != codeFmtLen {
		return ErrInvalidAuthToken
	}

	b := make([]byte, 0, codeTxtLen)
	b = append(b, text[:4]...)
	b = append(b, text[5:9]...)
	b = append(b, text[10:14]...)
	b = append(b, text[15:19]...)
	b = append(b, text[20:24]...)
	b = append(b, text[25:29]...)
	b = append(b, text[30:34]...)
	b = append(b, text[35:]...)

	_, err = encoder.Decode(t[:], b)
	return
}
