package main

import "testing"

//go:noinline
func doCallback(cb func(string)) {
	cb("foobar")
}

//go:noinline
func doInterface(iface foo) {
	iface.Do("foobar")
}

type foo interface {
	Do(string)
}

type bar struct{}

func (bar) Do(str string) {
	_ = str
}

//go:noinline
func cbfunc(s string) {
	_ = s
}

func BenchmarkCallback(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doCallback(func(s string) {
			_ = s
		})
	}
}

func BenchmarkCallback2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doCallback(cbfunc)
	}
}

func BenchmarkInterface(b *testing.B) {
	var f bar
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		doInterface(f)
	}
}

//go:noinline
func cb1(data []int) {
	_ = data
}

//go:noinline
func cb2(data ...int) {
	_ = data
}

func BenchmarkSliceCB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cb1([]int{1, 2, 3, 4, 5, 6, 7, 8})
	}
}

func BenchmarkSliceRest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cb2(1, 2, 3, 4, 5, 6, 7, 8)
	}
}
