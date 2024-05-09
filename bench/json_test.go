package main

import (
	"io"
	"testing"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/webmafia/fast"
)

func BenchmarkEncode_WriteVal(b *testing.B) {
	s := jsoniter.ConfigFastest.BorrowStream(io.Discard)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.WriteVal(User{
			ID:   1234,
			Name: "Hello there",
		})
	}
}

func BenchmarkEncode_Encoder(b *testing.B) {
	s := jsoniter.ConfigFastest.BorrowStream(io.Discard)
	typ := reflect2.TypeOf(User{})
	enc := jsoniter.ConfigFastest.EncoderOf(typ)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		enc.Encode(fast.Noescape(unsafe.Pointer(&User{
			ID:   1234,
			Name: "Hello there",
		})), s)
	}
}
