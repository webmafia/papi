package fastapi

import (
	"io"
	"log"
	"unsafe"

	"github.com/valyala/fasthttp"
)

var _ io.Closer = (*Params)(nil)

type Params struct {
	keys []string
	vals []string
}

// Implements io.Closer, so that fasthttp automatically resets this on release of RequestCtx.
func (p *Params) Close() error {
	p.Reset()
	return nil
}

func (p *Params) Reset() {
	p.keys = nil
	p.vals = p.vals[:0]
}

func (p *Params) Get(key string) (val string, ok bool) {
	for i := range p.keys {
		if p.keys[i] == key {
			return p.vals[i], true
		}
	}

	return
}

func RequestParams(c *fasthttp.RequestCtx) (params *Params) {
	var ok bool

	if params, ok = c.UserValue("params").(*Params); !ok {
		params = new(Params)
		c.SetUserValue("params", params)
	}

	log.Println(uintptr(unsafe.Pointer(c)), "got", uintptr(unsafe.Pointer(params)))

	return
}
