package token

import (
	"time"
	"unsafe"

	"github.com/webmafia/identifier"
)

type TokenView interface {
	Id() identifier.ID
	TimeCreated() time.Time
	Payload() [24]byte
}

type tokenView struct {
	ptr unsafe.Pointer
}

func (t tokenView) Id() identifier.ID {
	return *(*identifier.ID)(t.ptr)
}

func (t tokenView) TimeCreated() time.Time {
	return t.Id().Time()
}

func (t tokenView) Payload() [24]byte {
	return *(*[24]byte)(unsafe.Add(t.ptr, 8))
}

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}
