package json

import (
	"io"
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

var pool = NewJSONPool(jsoniter.ConfigFastest)

//go:inline
func AcquireIterator(r io.Reader) *jsoniter.Iterator {
	return pool.AcquireIterator(r)
}

//go:inline
func ReleaseIterator(iter *jsoniter.Iterator) {
	pool.ReleaseIterator(iter)
}

//go:inline
func AcquireStream(w io.Writer) *jsoniter.Stream {
	return pool.AcquireStream(w)
}

//go:inline
func ReleaseStream(s *jsoniter.Stream) {
	pool.ReleaseStream(s)
}

//go:inline
func DecoderOf(typ reflect.Type) jsoniter.ValDecoder {
	return pool.DecoderOf(typ)
}

//go:inline
func EncoderOf(typ reflect.Type) jsoniter.ValEncoder {
	return pool.EncoderOf(typ)
}
