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
	route *route[U]
}

type route[U any] struct {
	params []string
	cb     func(ctx *Ctx[U]) error
}

func (r *Router[U]) Clear() {
	r.tree.nodes = r.tree.nodes[:0]
}

func (r *Router[U]) Add(method string, path string) *route[U] {
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
	n.route = &route[U]{
		params: params,
	}

	return n.route
}

func (r *Router[U]) add(n *node[U], part []byte, params *[]string) *node[U] {
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

func (r *Router[U]) LookupString(method string, path string, paramVals *[]string) (cb func(ctx *Ctx[U]) error, params []string, ok bool) {
	return r.Lookup(fast.StringToBytes(method), fast.StringToBytes(path), paramVals)
}

func (r *Router[U]) Lookup(method []byte, p []byte, paramVals *[]string) (cb func(ctx *Ctx[U]) error, params []string, ok bool) {
	if p[0] == '/' {
		p = p[1:]
	}

	n := r.lookup(&r.tree, method, paramVals)

	if n == nil {
		return
	}

	for {
		idx := bytes.IndexByte(p, '/')

		if idx < 0 {
			break
		}

		n = r.lookup(n, p[:idx], paramVals)

		if n == nil {
			return
		}

		p = p[idx+1:]
	}

	n = r.lookup(n, p, paramVals)

	if n == nil {
		return
	}

	return n.route.cb, n.route.params, true
}

func (r *Router[U]) lookup(n *node[U], part []byte, paramVals *[]string) *node[U] {
	if len(part) == 0 {
		return n
	}

	for i := range n.nodes {
		if n.nodes[i].value[0] == '{' {
			*paramVals = append(*paramVals, fast.BytesToString(part))
			return n.nodes[i]
		}

		if bytes.Equal(n.nodes[i].value, part) {
			return n.nodes[i]
		}
	}

	return nil
}
