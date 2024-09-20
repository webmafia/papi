package internal

import (
	"reflect"
	"unsafe"
)

// Checks whether type T implements interface I.
//
//go:inline
func ImplementsInterface[T any, I any]() (ok bool) {
	_, ok = any((*T)(nil)).(I)
	return
}

// Gets a reflect.Type of T.
//
//go:inline
func ReflectType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// Panics if err != nil.
//
//go:inline
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func ptrToInterface[T any](dst *T, src unsafe.Pointer) {
	typ := reflect.TypeOf((*T)(nil)).Elem()

	// Create a reflect.Value from the unsafe.Pointer
	val := reflect.NewAt(typ, src).Elem()

	// Handle both pointer and value types
	if val.Kind() == reflect.Ptr {
		*dst = val.Interface().(T) // If it's already a pointer, assign it directly
	} else {
		*dst = *val.Addr().Interface().(*T) // If it's not a pointer, take the address and assign
	}
}
