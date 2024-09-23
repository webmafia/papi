package fastapi

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
	route *route
}

type route struct {
	params  []string
	handler func(c *fasthttp.RequestCtx) error
}

func (r *Router) Clear() {
	r.tree.nodes = r.tree.nodes[:0]
}

func (r *Router) Add(method string, path string) *route {
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
	n.route = &route{
		params: params,
	}

	return n.route
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

func (r *Router) LookupString(method string, path string, paramVals *[]string) (cb func(c *fasthttp.RequestCtx) error, params []string, ok bool) {
	return r.Lookup(fast.StringToBytes(method), fast.StringToBytes(path), paramVals)
}

func (r *Router) Lookup(method []byte, p []byte, paramVals *[]string) (cb func(c *fasthttp.RequestCtx) error, params []string, ok bool) {
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

	return n.route.handler, n.route.params, true
}

func (r *Router) lookup(n *node, part []byte, paramVals *[]string) *node {
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
