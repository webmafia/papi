package route

import (
	"bytes"
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/fast"
)

type Router struct {
	tree      node
	paramPool sync.Pool
}

type node struct {
	nodes []*node
	value []byte
	route *Route
}

type Route struct {
	Params  []string
	Handler func(c *fasthttp.RequestCtx) error
}

func (r *Router) Clear() {
	r.tree.nodes = r.tree.nodes[:0]
}

func (r *Router) Add(method string, path string, fn func(*Route) error) error {
	var params []string
	p := fast.StringToBytes(path)
	n := r.add(&r.tree, fast.StringToBytes(method), &params)

	for {
		idx := bytes.IndexByte(p, '/')

		if idx < 0 {
			break
		}

		n = r.add(n, p[:idx], &params)
		p = p[idx+1:]
	}

	n = r.add(n, p, &params)
	n.route = &Route{
		Params: params,
	}

	return fn(n.route)
}

func (r *Router) add(n *node, part []byte, params *[]string) *node {
	if len(part) == 0 {
		return n
	}

	if part[0] == '{' {
		*params = append(*params, fast.BytesToString(part[1:len(part)-1]))
	}

	for i := range n.nodes {
		if bytes.Equal(n.nodes[i].value, part) {
			return n.nodes[i]
		}
	}

	nn := &node{
		value: part,
	}

	// Routes with params go last
	if len(n.nodes) == 0 || nn.value[0] == '{' {
		n.nodes = append(n.nodes, nn)
	} else {
		n.nodes = append([]*node{nn}, n.nodes...)
	}

	return nn
}

func (r *Router) LookupString(method string, path string) (cb func(c *fasthttp.RequestCtx) error, params *Params, ok bool) {
	return r.Lookup(fast.StringToBytes(method), fast.StringToBytes(path))
}

func (r *Router) Lookup(method []byte, p []byte) (cb func(c *fasthttp.RequestCtx) error, params *Params, ok bool) {
	params = r.acquireParams()

	if p[0] == '/' {
		p = p[1:]
	}

	n := r.lookup(&r.tree, method, params)

	if n == nil {
		return
	}

	for {
		idx := bytes.IndexByte(p, '/')

		if idx < 0 {
			break
		}

		n = r.lookup(n, p[:idx], params)

		if n == nil {
			return
		}

		p = p[idx+1:]
	}

	n = r.lookup(n, p, params)

	if n == nil || n.route == nil {
		return
	}

	params.keys = n.route.Params

	return n.route.Handler, params, true
}

func (r *Router) lookup(n *node, part []byte, params *Params) *node {
	if len(part) == 0 {
		return n
	}

	for i := range n.nodes {
		if n.nodes[i].value[0] == '{' {
			params.addValue(part)
			return n.nodes[i]
		}

		if bytes.Equal(n.nodes[i].value, part) {
			return n.nodes[i]
		}
	}

	return nil
}

func (r *Router) acquireParams() (p *Params) {
	var ok bool

	if p, ok = r.paramPool.Get().(*Params); !ok {
		p = newParams(4)
	}

	return
}

func (r *Router) ReleaseParams(p *Params) {
	if p != nil {
		p.Reset()
		r.paramPool.Put(p)
	}
}
