package registry

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/webbmaffian/papi/pool/json"
	"github.com/webbmaffian/papi/registry/scanner"
)

type Registry struct {
	req map[reflect.Type]RequestType
	typ map[reflect.Type]ParamType
	tag map[reflect.Type]ParamDecoder
	def *requestScanner
	mu  sync.RWMutex
}

func NewRegistry(json *json.Pool) (r *Registry, err error) {
	r = &Registry{
		req: make(map[reflect.Type]RequestType),
		typ: make(map[reflect.Type]ParamType),
		tag: make(map[reflect.Type]ParamDecoder),
	}

	r.def = &requestScanner{
		reg:  r,
		json: json,
	}

	return
}

//go:inline
func (s *Registry) JSON() *json.Pool {
	return s.def.json
}

func (s *Registry) CreateRequestDecoder(typ reflect.Type, tags reflect.StructTag, paramKeys []string, fallback ...bool) (scan RequestDecoder, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if creator, ok := s.req[typ]; ok {
		return creator.CreateRequestDecoder(tags, paramKeys)
	}

	if len(fallback) > 0 && fallback[0] && s.def != nil {
		return s.def.CreateRequestDecoder(typ, tags, paramKeys)
	}

	return nil, errors.New("no request decoder could be found nor created")
}

func (s *Registry) CreateResponseEncoder(typ reflect.Type, tags reflect.StructTag, paramKeys []string, handler ResponseEncoder, fallback ...bool) (scan ResponseEncoder, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if creator, ok := s.req[typ]; ok {
		return creator.CreateResponseEncoder(s, tags, paramKeys, handler)
	}

	if len(fallback) > 0 && fallback[0] && s.def != nil {
		return s.def.CreateResponseEncoder(typ, tags, paramKeys, handler)
	}

	return nil, errors.New("no response encoder could be found nor created")
}

func (s *Registry) RegisterType(typs ...Type) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, typ := range typs {
		var reg bool

		if typ, ok := typ.(ParamType); ok {
			s.typ[typ.Type()] = typ
			reg = true
		}

		if typ, ok := typ.(RequestType); ok {
			s.req[typ.Type()] = typ
			reg = true
		}

		if !reg {
			return fmt.Errorf("%T is neither a ParamType nor RequestType", typ)
		}
	}

	return
}

func (s *Registry) CreateParamDecoder(typ reflect.Type, tags reflect.StructTag) (scan ParamDecoder, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var createScanner scanner.CreateValueScanner

	createScanner = func(typ reflect.Type, createElemScanner scanner.CreateValueScanner) (scan ParamDecoder, err error) {
		if creator, ok := s.typ[typ]; ok {
			return creator.CreateParamDecoder(tags)
		}

		return scanner.CreateCustomScanner(typ, createScanner)
	}

	return createScanner(typ, createScanner)
}
