package papi

import (
	"io"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func BenchmarkList_Write(b *testing.B) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var l List[User]

	s := jsoniter.ConfigFastest.BorrowStream(io.Discard)
	l.setStream(s)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Write(&User{
			ID:   1234,
			Name: "Hello there",
		})
	}
}
