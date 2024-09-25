package scanner

import (
	"errors"
	"reflect"
	"sync"

	"github.com/webmafia/fastapi/scanner/value"
)

type Registry struct {
	req    sync.Map
	val    sync.Map
	def    RequestScannerCreator
	frozen bool
}

func NewRegistry(init ...func(*Registry)) (r *Registry) {
	r = &Registry{}

	if len(init) > 0 {
		init[0](r)
	}

	r.frozen = true
	return r
}

func (s *Registry) RegisterDefaultRequestScanner(creator RequestScannerCreator) {
	if s.frozen {
		return
	}

	s.def = creator
}

func (s *Registry) RegisterRequestScanner(typ reflect.Type, creator RequestScannerCreator) {
	if creator == nil {
		s.req.Delete(typ)
	} else {
		s.req.Store(typ, creator)
	}
}

func (s *Registry) CreateRequestScanner(typ reflect.Type, tags reflect.StructTag, paramKeys []string, fallback ...bool) (scan RequestScanner, err error) {
	if v, ok := s.req.Load(typ); ok {
		if creator, ok := v.(RequestScannerCreator); ok {
			return creator.CreateScanner(typ, tags, paramKeys)
		}

		return nil, errors.New("invalid request scanner creator - this should not be possible")
	}

	if len(fallback) > 0 && fallback[0] && s.def != nil {
		return s.def.CreateScanner(typ, tags, paramKeys)
	}

	return nil, errors.New("no scanner could be found nor created")
}

func (s *Registry) RegisterValueScanner(typ reflect.Type, create CreateValueScanner) {
	if create == nil {
		s.req.Delete(typ)
	} else {
		s.req.Store(typ, create)
	}
}

func (s *Registry) CreateValueScanner(typ reflect.Type, tags reflect.StructTag) (scan value.ValueScanner, err error) {
	var createScanner value.CreateValueScanner

	createScanner = func(typ reflect.Type, createElemScanner value.CreateValueScanner) (scan value.ValueScanner, err error) {
		if v, ok := s.val.Load(typ); ok {
			if create, ok := v.(CreateValueScanner); ok {
				return create(tags)
			}

			return nil, errors.New("invalid value scanner creator - this should not be possible")
		}

		return value.CreateCustomScanner(typ, createScanner)
	}

	return createScanner(typ, createScanner)
}
