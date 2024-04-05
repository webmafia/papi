package router

import (
	"bytes"
	"unsafe"

	"github.com/webmafia/fast"
)

type Router struct {
	tree node
}

func New() *Router {
	return &Router{}
}

type node struct {
	nodes []*node
	value []byte
	ptr   unsafe.Pointer
}

func (r *Router) Clear() {
	r.tree.nodes = r.tree.nodes[:0]
}

func (r *Router) Add(method string, path string, ptr unsafe.Pointer) {
	p := fast.StringToBytes(path)
	n := r.add(&r.tree, fast.StringToBytes(method))

	for {
		idx := bytes.IndexByte(p, '/')

		if idx < 0 {
			break
		}

		n = r.add(n, p[:idx])
		p = p[idx+1:]
	}

	n = r.add(n, p)
	n.ptr = ptr
}

func (r *Router) add(n *node, part []byte) *node {
	if len(part) == 0 {
		return n
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

func (r *Router) LookupString(method string, path string, params *Params) (ptr unsafe.Pointer) {
	return r.Lookup(fast.StringToBytes(method), fast.StringToBytes(path), params)
}

func (r *Router) Lookup(method []byte, p []byte, params *Params) (ptr unsafe.Pointer) {
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

	return n.ptr
}

func (r *Router) lookup(n *node, part []byte, params *Params) *node {
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
