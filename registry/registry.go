package registry

import (
	"errors"
	"reflect"
	"sync"

	"github.com/webmafia/fastapi/openapi"
	"github.com/webmafia/fastapi/registry/value"
)

type Registry struct {
	req     map[reflect.Type]RequestScannerCreator
	val     map[reflect.Type]CreateValueScanner
	schemas map[reflect.Type]*openapi.Schema
	def     RequestScannerCreator
	mu      sync.RWMutex
}

func NewRegistry(def ...func(*Registry) RequestScannerCreator) (r *Registry) {
	r = &Registry{
		req: make(map[reflect.Type]RequestScannerCreator),
		val: make(map[reflect.Type]CreateValueScanner),
	}

	if len(def) > 0 {
		r.def = def[0](r)
	}

	return r
}

func (s *Registry) RegisterRequestScanner(typ reflect.Type, creator RequestScannerCreator) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if creator == nil {
		delete(s.req, typ)
	} else {
		s.req[typ] = creator
	}
}

func (s *Registry) CreateRequestScanner(typ reflect.Type, tags reflect.StructTag, paramKeys []string, fallback ...bool) (scan RequestScanner, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if creator, ok := s.req[typ]; ok {
		return creator.CreateScanner(typ, tags, paramKeys)
	}

	if len(fallback) > 0 && fallback[0] && s.def != nil {
		return s.def.CreateScanner(typ, tags, paramKeys)
	}

	return nil, errors.New("no scanner could be found nor created")
}

func (s *Registry) DescribeOperation(op *openapi.Operation, in, out reflect.Type) (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Input
	if creator, ok := s.req[in]; ok {
		err = creator.Describe(op, in)
	} else if s.def != nil {
		err = s.def.Describe(op, in)
	} else {
		err = errors.New("no input descriptor could be found nor created")
	}

	if err != nil {
		return
	}

	// Output
	if op.Response, err = s.Schema(out); err != nil {
		return
	}

	return
}

func (s *Registry) RegisterValueScanner(typ reflect.Type, create CreateValueScanner) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if create == nil {
		delete(s.val, typ)
	} else {
		s.val[typ] = create
	}
}

func (s *Registry) CreateValueScanner(typ reflect.Type, tags reflect.StructTag) (scan value.ValueScanner, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var createScanner value.CreateValueScanner

	createScanner = func(typ reflect.Type, createElemScanner value.CreateValueScanner) (scan value.ValueScanner, err error) {
		if create, ok := s.val[typ]; ok {
			return create(tags)
		}

		return value.CreateCustomScanner(typ, createScanner)
	}

	return createScanner(typ, createScanner)
}
