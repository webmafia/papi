package router

import "github.com/webmafia/fast"

type Params []kv

func (p *Params) Reset() {
	*p = (*p)[:0]
}

func (p *Params) add(key, value []byte) {
	*p = append(*p, kv{
		key:   fast.BytesToString(key),
		value: fast.BytesToString(value),
	})
}

func (p Params) Get(key string) (val string, ok bool) {
	for i := range p {
		if p[i].key == key {
			return p[i].value, true
		}
	}

	return
}

type kv struct {
	key   string
	value string
}
