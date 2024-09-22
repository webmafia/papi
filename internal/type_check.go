package internal

import (
	"reflect"
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
