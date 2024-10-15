package main

import (
	"iter"
	"testing"
)

type bar struct {
	name int
}

type foo[T any] struct {
	list iter.Seq[*T]
}

func BenchmarkCallback(b *testing.B) {
	var f foo[bar]

	b.ResetTimer()

	for range b.N {
		callback(&f)

		for v := range f.list {
			_ = v
		}
	}
}

func BenchmarkReturn(b *testing.B) {
	for range b.N {
		for v := range returnIter().list {
			_ = v
		}
	}
}

func BenchmarkReturnPlain(b *testing.B) {
	for range b.N {
		for v := range returnPlainIter() {
			_ = v
		}
	}
}

func callback(f *foo[bar]) {
	f.list = func(yield func(*bar) bool) {
		v := bar{name: 123}
		yield(&v)
	}
}

func returnIter() foo[bar] {
	return foo[bar]{
		list: func(yield func(*bar) bool) {
			v := bar{name: 123}
			yield(&v)
		},
	}
}

func returnPlainIter() iter.Seq[*bar] {
	return func(yield func(*bar) bool) {
		v := bar{name: 123}
		yield(&v)
	}
}
