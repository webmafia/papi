package registry

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/webbmaffian/papi/pool/json"
	"github.com/webbmaffian/papi/registry/scanner"
	"github.com/webbmaffian/papi/registry/types"
)

type Registry struct {
	req map[reflect.Type]types.RequestType
	typ map[reflect.Type]types.ParamType
	tag map[reflect.Type]types.ParamDecoder
	def *requestScanner
	mu  sync.RWMutex
}

func NewRegistry(json *json.Pool) (r *Registry, err error) {
	r = &Registry{
		req: make(map[reflect.Type]types.RequestType),
		typ: make(map[reflect.Type]types.ParamType),
		tag: make(map[reflect.Type]types.ParamDecoder),
	}

	r.def = &requestScanner{
		reg:  r,
		json: json,
	}

	if err = r.registerCommonTypes(); err != nil {
		return
	}

	return
}

func (s *Registry) CreateRequestDecoder(typ reflect.Type, tags reflect.StructTag, paramKeys []string, fallback ...bool) (scan types.RequestDecoder, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if creator, ok := s.req[typ]; ok {
		return creator.CreateRequestDecoder(tags, paramKeys)
	}

	if len(fallback) > 0 && fallback[0] && s.def != nil {
		return s.def.CreateRequestDecoder(typ, tags, paramKeys)
	}

	return nil, errors.New("no scanner could be found nor created")
}

func (s *Registry) registerCommonTypes() (err error) {
	return s.RegisterType(
		types.TimeType(),
	)
}

func (s *Registry) RegisterType(typs ...types.Type) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, typ := range typs {
		var reg bool

		if typ, ok := typ.(types.ParamType); ok {
			s.typ[typ.Type()] = typ
			reg = true
		}

		if typ, ok := typ.(types.RequestType); ok {
			s.req[typ.Type()] = typ
			reg = true
		}

		if !reg {
			return fmt.Errorf("%T is neither a ParamType nor RequestType", typ)
		}
	}

	return
}

func (s *Registry) CreateParamDecoder(typ reflect.Type, tags reflect.StructTag) (scan types.ParamDecoder, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var createScanner scanner.CreateValueScanner

	createScanner = func(typ reflect.Type, createElemScanner scanner.CreateValueScanner) (scan types.ParamDecoder, err error) {
		if creator, ok := s.typ[typ]; ok {
			return creator.CreateParamDecoder(tags)
		}

		return scanner.CreateCustomScanner(typ, createScanner)
	}

	return createScanner(typ, createScanner)
}
