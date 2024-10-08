package registry

import (
	"errors"
	"reflect"
	"sync"

	"github.com/webmafia/fastapi/registry/scanner"
	"github.com/webmafia/fastapi/registry/types"
)

type Registry struct {
	req map[reflect.Type]RequestScannerCreator
	typ map[reflect.Type]types.Type
	def RequestScannerCreator
	mu  sync.RWMutex
}

func NewRegistry(def ...func(*Registry) RequestScannerCreator) (r *Registry) {
	r = &Registry{
		req: make(map[reflect.Type]RequestScannerCreator),
		typ: make(map[reflect.Type]types.Type),
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

func (s *Registry) RegisterCommonTypes() {
	s.RegisterType(
		types.TimeType(),
	)
}

func (s *Registry) RegisterType(types ...types.Type) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, typ := range types {
		if typ != nil {
			s.typ[typ.Type()] = typ
		}
	}
}

func (s *Registry) CreateValueScanner(typ reflect.Type, tags reflect.StructTag) (scan scanner.Scanner, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var createScanner scanner.CreateValueScanner

	createScanner = func(typ reflect.Type, createElemScanner scanner.CreateValueScanner) (scan scanner.Scanner, err error) {
		if creator, ok := s.typ[typ]; ok {
			return creator.CreateScanner(tags)
		}

		return scanner.CreateCustomScanner(typ, createScanner)
	}

	return createScanner(typ, createScanner)
}
