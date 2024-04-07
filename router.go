package fastapi

import (
	"bytes"

	"github.com/webmafia/fast"
)

type Router[U any] struct {
	tree node[U]
}

type node[U any] struct {
	nodes []*node[U]
	value []byte
	cb    func(ctx *Ctx[U]) error
}

func (r *Router[U]) Clear() {
	r.tree.nodes = r.tree.nodes[:0]
}

func (r *Router[U]) Add(method string, path string, cb func(ctx *Ctx[U]) error, paramsCb ...func(string)) {
	p := fast.StringToBytes(path)
	n := r.add(&r.tree, fast.StringToBytes(method))

	for {
		idx := bytes.IndexByte(p, '/')

		if idx < 0 {
			break
		}

		n = r.add(n, p[:idx], paramsCb...)
		p = p[idx+1:]
	}

	n = r.add(n, p, paramsCb...)
	n.cb = cb
}

func (r *Router[U]) add(n *node[U], part []byte, paramsCb ...func(string)) *node[U] {
	if len(part) == 0 {
		return n
	}

	if part[0] == '{' && len(paramsCb) > 0 {
		paramsCb[0](fast.BytesToString(part[1 : len(part)-1]))
	}

	for i := range n.nodes {
		if bytes.Equal(n.nodes[i].value, part) {
			return n.nodes[i]
		}
	}

	nn := &node[U]{
		value: part,
	}

	// Routes with params go last
	if len(n.nodes) == 0 || nn.value[0] == '{' {
		n.nodes = append(n.nodes, nn)
	} else {
		n.nodes = append([]*node[U]{nn}, n.nodes...)
	}

	return nn
}

func (r *Router[U]) LookupString(method string, path string, params *Params) func(ctx *Ctx[U]) error {
	return r.Lookup(fast.StringToBytes(method), fast.StringToBytes(path), params)
}

func (r *Router[U]) Lookup(method []byte, p []byte, params *Params) (cb func(ctx *Ctx[U]) error) {
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

	if n == nil {
		return
	}

	return n.cb
}

func (r *Router[U]) lookup(n *node[U], part []byte, params *Params) *node[U] {
	if len(part) == 0 {
		return n
	}

	for i := range n.nodes {
		if n.nodes[i].value[0] == '{' {
			params.add(n.nodes[i].value[1:len(n.nodes[i].value)-1], part)
			return n.nodes[i]
		}

		if bytes.Equal(n.nodes[i].value, part) {
			return n.nodes[i]
		}
	}

	return nil
}
