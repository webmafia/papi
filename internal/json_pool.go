package internal

import (
	"io"
	"os"
	"reflect"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/webmafia/fast"
)

type JSONPool struct {
	api      jsoniter.API
	iterPool sync.Pool
}

func NewJSONPool(api jsoniter.API, bufSize ...int) *JSONPool {
	var size int

	if len(bufSize) > 0 {
		size = bufSize[0]
	}

	if size <= 0 {
		size = os.Getpagesize()
	}

	return &JSONPool{
		api: api,
		iterPool: sync.Pool{
			New: func() any {
				return jsoniter.
					NewIterator(jsoniter.ConfigFastest).
					ResetBytes(fast.MakeNoZero(size))
			},
		},
	}
}

func (p *JSONPool) AcquireIterator(r io.Reader) *jsoniter.Iterator {
	return p.iterPool.Get().(*jsoniter.Iterator).Reset(r)
}

func (p *JSONPool) ReleaseIterator(iter *jsoniter.Iterator) {
	iter.Error = nil
	p.iterPool.Put(iter.Reset(nil))
}

//go:inline
func (p *JSONPool) AcquireStream(w io.Writer) *jsoniter.Stream {
	return p.api.BorrowStream(w)
}

//go:inline
func (p *JSONPool) ReleaseStream(s *jsoniter.Stream) {
	s.Error = nil
	p.api.ReturnStream(s)
}

func (p *JSONPool) DecoderOf(typ reflect.Type) jsoniter.ValDecoder {
	return p.api.DecoderOf(reflect2.Type2(reflect.PointerTo(typ)))
}

func (p *JSONPool) EncoderOf(typ reflect.Type) jsoniter.ValEncoder {
	return p.api.EncoderOf(reflect2.Type2(reflect.PointerTo(typ)))
}
