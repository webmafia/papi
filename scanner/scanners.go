package scanner

import (
	"errors"
	"reflect"
	"sync"

	"github.com/webmafia/fastapi/scanner/value"
)

type ScannerCreator struct {
	req sync.Map
	val sync.Map
}

func NewScannerCreator() *ScannerCreator {
	return &ScannerCreator{}
}

func (s *ScannerCreator) RegisterRequestScanner(typ reflect.Type, creator RequestScannerCreator) {
	if creator == nil {
		s.req.Delete(typ)
	} else {
		s.req.Store(typ, creator)
	}
}

func (s *ScannerCreator) CreateRequestScanner(typ reflect.Type, tags reflect.StructTag, paramKeys []string) (scan RequestScanner, err error) {
	if v, ok := s.req.Load(typ); ok {
		if creator, ok := v.(RequestScannerCreator); ok {
			return creator.CreateScanner(paramKeys, tags)
		}

		return nil, errors.New("invalid request scanner creator - this should not be possible")
	}

	return nil, errors.New("no scanner could be found nor created")
}

func (s *ScannerCreator) RegisterValueScanner(typ reflect.Type, create CreateValueScanner) {
	if create == nil {
		s.req.Delete(typ)
	} else {
		s.req.Store(typ, create)
	}
}

func (s *ScannerCreator) CreateValueScanner(typ reflect.Type, tags reflect.StructTag) (scan value.ValueScanner, err error) {
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
