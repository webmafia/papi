package jsonpool

import (
	"io"
	"os"
	"reflect"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/webmafia/fast"
)

var (
	pageSize = os.Getpagesize()
	iterPool = sync.Pool{
		New: func() any {
			return jsoniter.
				NewIterator(jsoniter.ConfigFastest).
				ResetBytes(fast.MakeNoZero(pageSize))
		},
	}
)

func AcquireIterator(r io.Reader) *jsoniter.Iterator {
	return iterPool.Get().(*jsoniter.Iterator).Reset(r)
}

func ReleaseIterator(iter *jsoniter.Iterator) {
	iterPool.Put(iter.Reset(nil))
}

//go:inline
func AcquireStream(w io.Writer) *jsoniter.Stream {
	return jsoniter.ConfigFastest.BorrowStream(w)
}

//go:inline
func ReleaseStream(s *jsoniter.Stream) {
	jsoniter.ConfigFastest.ReturnStream(s)
}

func DecoderOf(typ reflect.Type) jsoniter.ValDecoder {
	return jsoniter.ConfigFastest.DecoderOf(reflect2.Type2(reflect.PointerTo(typ)))
}
