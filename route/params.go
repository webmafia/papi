package route

import "github.com/webmafia/fast"

var NilParams = &Params{}

type Params struct {
	keys []string
	vals []string
}

func newParams(c int) *Params {
	return &Params{
		vals: make([]string, 0, c),
	}
}

func (p *Params) addValue(val []byte) {
	p.vals = append(p.vals, fast.BytesToString(val))
}

func (p *Params) Reset() {
	p.keys = nil
	p.vals = p.vals[:0]
}

func (p *Params) Value(idx int) string {
	if idx < len(p.vals) {
		return p.vals[idx]
	}

	return ""
}

func (p *Params) Get(key string) (val string, ok bool) {
	for i := range p.keys {
		if p.keys[i] == key {
			return p.vals[i], true
		}
	}

	return
}

func (p *Params) Valid() bool {
	return len(p.keys) == len(p.vals)
}
