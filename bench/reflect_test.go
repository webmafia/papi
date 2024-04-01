package main

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"unsafe"

	"github.com/webmafia/fastapi"
)

type foobar struct{}

//go:noinline
func (f foobar) ListUsers(api *fastapi.API[User]) {
	// fmt.Println("Tadaaaaa")
}

func ExampleReflect_Call() {
	f := foobar{}
	v := reflect.ValueOf(f).Method(0).Interface()
	cb, ok := v.(func(api *fastapi.API[User]))

	if !ok {
		panic("nope")
	}

	cb(nil)

	// t := reflect.TypeOf(foobar{})
	// ptr := unsafe.Add(unsafe.Pointer(&foobar{}), t.Method(0).Type.Align())
	// cb := *(*func(api *fastapi.API[User]))(ptr)

	// cb(nil)

	// Output: Trudelutt
}

func ExampleReflect_MethodName() {
	f := foobar{}
	p := reflect.ValueOf(f.ListUsers).Pointer()
	fmt.Printf("%#v\n", runtime.FuncForPC(p).Name())

	// Output: Trudelutt
}

type MethodValue struct {
	Receiver uintptr
	FuncPtr  uintptr
}

func ExampleReflect_MethodName_Unsafe() {
	f := foobar{}
	cb := f.ListUsers
	p := *(*MethodValue)(unsafe.Pointer(&cb))
	fmt.Printf("%#v\n", runtime.FuncForPC(p.FuncPtr).Name())

	// t := reflect.TypeOf(foobar{})
	// ptr := unsafe.Add(unsafe.Pointer(&foobar{}), t.Method(0).Type.Align())
	// cb := *(*func(api *fastapi.API[User]))(ptr)

	// cb(nil)

	// Output: Trudelutt
}

func BenchmarkReflect_Name(b *testing.B) {
	f := foobar{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = reflect.ValueOf(f.ListUsers).Pointer()
	}
}

func BenchmarkReflect_Method(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := foobar{}
		v := reflect.ValueOf(f).Method(0).Interface()
		cb, ok := v.(func(api *fastapi.API[User]))

		if !ok {
			b.Fatal("nope")
		}

		_ = cb
	}
}

func BenchmarkReflect_MethodCall(b *testing.B) {
	f := foobar{}
	v := reflect.ValueOf(f).Method(0).Interface()
	cb, ok := v.(func(api *fastapi.API[User]))

	if !ok {
		b.Fatal("nope")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cb(nil)
	}
}

func BenchmarkReflect_MethodCall_Unsafe(b *testing.B) {
	f := foobar{}
	v := reflect.ValueOf(f).Method(0).UnsafePointer()
	cb := *(*func(api *fastapi.API[User]))(v)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cb(nil)
	}
}
