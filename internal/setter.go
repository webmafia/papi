package internal

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/modern-go/reflect2"
)

//go:linkname typedmemclr reflect.typedmemclr
func typedmemclr(rtype unsafe.Pointer, ptr unsafe.Pointer)

type Setter struct {
	typ reflect2.Type
	ptr unsafe.Pointer
}

// Created a new setter that assumes ptr will live. It is NOT safe to let the setter
// outlive ptr.
func NewSetter(typ reflect.Type, ptr unsafe.Pointer) Setter {
	return Setter{
		typ: reflect2.Type2(typ),
		ptr: reflect2.NoEscape(ptr),
	}
}

// Returns the type of the setter.
func (s Setter) Type() reflect.Type {
	return s.typ.Type1()
}

// Sets the setter to a specific value. The value MUST be wrapped in a pointer - any
// non-pointers will return an error. Passing a nil value will zero the destination.
func (s Setter) Set(v any) (err error) {
	if v == nil {
		// Zero the destination value of type s.typ at s.ptr.
		typedmemclr(unsafe.Pointer(s.typ.RType()), s.ptr)
		return
	}

	typ := reflect.TypeOf(v)

	if typ.Kind() != reflect.Pointer {
		return errors.New("setter: value must be a pointer")
	}

	typ = typ.Elem()

	if typ != s.typ.Type1() {
		return fmt.Errorf("setter: mismatching type; expected %s, got %s", s.typ, typ)
	}

	s.typ.UnsafeSet(s.ptr, reflect2.NoEscape(reflect2.PtrOf(v)))
	return
}
