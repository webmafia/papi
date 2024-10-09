package json

import (
	"io"
	"os"
	"reflect"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/webmafia/fast"
)

type Pool struct {
	api      jsoniter.API
	iterPool sync.Pool
}

func NewPool(api jsoniter.API, bufSize ...int) *Pool {
	var size int

	if len(bufSize) > 0 {
		size = bufSize[0]
	}

	if size <= 0 {
		size = os.Getpagesize()
	}

	return &Pool{
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

func (p *Pool) AcquireIterator(r io.Reader) *jsoniter.Iterator {
	return p.iterPool.Get().(*jsoniter.Iterator).Reset(r)
}

func (p *Pool) ReleaseIterator(iter *jsoniter.Iterator) {
	iter.Error = nil
	p.iterPool.Put(iter.Reset(nil))
}

//go:inline
func (p *Pool) AcquireStream(w io.Writer) *jsoniter.Stream {
	return p.api.BorrowStream(w)
}

//go:inline
func (p *Pool) ReleaseStream(s *jsoniter.Stream) {
	s.Error = nil
	p.api.ReturnStream(s)
}

func (p *Pool) DecoderOf(typ reflect.Type) jsoniter.ValDecoder {
	return p.api.DecoderOf(reflect2.Type2(reflect.PointerTo(typ)))
}
