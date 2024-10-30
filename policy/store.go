package policy

import (
	"bytes"
	"errors"
	"fmt"
	"iter"
	"reflect"
	"sync"
	"unsafe"

	"github.com/webmafia/papi/internal"
)

type Store struct {
	store map[PolicyKey]unsafe.Pointer
	types map[accessKey]reflect.Type
	json  *internal.JSONPool
	mu    sync.RWMutex
}

func NewStore(json *internal.JSONPool) *Store {
	return &Store{
		types: make(map[accessKey]reflect.Type),
		json:  json,
	}
}

func (s *Store) Register(action, resource string, typ reflect.Type) (err error) {
	if action == "" {
		return errors.New("missing 'action' for policy")
	}

	if resource == "" {
		return errors.New("missing 'resource' for policy")
	}

	key := accessKey{action: action, resource: resource}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.store != nil {
		return errors.New("store has already been frozen and can't accept any new registration")
	}

	if regTyp, ok := s.types[key]; ok {
		if regTyp != typ {
			err = fmt.Errorf("can't register %s as policy for %s:%s, as %s already is registered", typ, action, resource, regTyp)
		}

		return
	}

	s.types[key] = typ
	return
}

func (s *Store) Freeze() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.freeze()
}

func (s *Store) freeze() {
	if s.store == nil {
		s.store = make(map[PolicyKey]unsafe.Pointer)
	}
}

func (s *Store) Add(role, action, resource string, condJson []byte) (err error) {
	typ, ok := s._GetType(action, resource)

	if !ok {
		return fmt.Errorf("no policy for %s:%s found", action, resource)
	}

	var read bytes.Reader
	read.Reset(condJson)

	dec := s.json.DecoderOf(typ)
	ptr := reflect.New(typ).UnsafePointer()
	iter := s.json.AcquireIterator(&read)
	defer s.json.ReleaseIterator(iter)

	dec.Decode(ptr, iter)

	if iter.Error != nil {
		return iter.Error
	}

	s._Set(role, action, resource, ptr)
	return
}

func (s *Store) BatchAdd(cb func(add func(role, action, resource string, condJson []byte) error) error) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.freeze()

	var read bytes.Reader
	iter := s.json.AcquireIterator(nil)
	defer s.json.ReleaseIterator(iter)

	add := func(role, action, resource string, condJson []byte) error {
		read.Reset(condJson)
		iter.Reset(&read)

		typ, ok := s.getType(action, resource)

		if !ok {
			return fmt.Errorf("no policy for %s:%s found", action, resource)
		}

		dec := s.json.DecoderOf(typ)
		ptr := reflect.New(typ).UnsafePointer()

		dec.Decode(ptr, iter)

		if iter.Error != nil {
			return iter.Error
		}

		s.set(role, action, resource, ptr)
		return nil
	}

	return cb(add)
}

func (s *Store) _GetType(action, resource string) (typ reflect.Type, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.getType(action, resource)
}

func (s *Store) getType(action, resource string) (typ reflect.Type, ok bool) {
	typ, ok = s.types[accessKey{action: action, resource: resource}]
	return
}

func (s *Store) _Set(role, action, resource string, condJson unsafe.Pointer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.freeze()
	s.set(role, action, resource, condJson)
}

func (s *Store) set(role, action, resource string, condJson unsafe.Pointer) {
	s.store[PolicyKey{Role: role, Action: action, Resource: resource}] = condJson
}

func (s *Store) Iterate() iter.Seq2[PolicyKey, []byte] {
	return func(yield func(PolicyKey, []byte) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		dumped := make(map[accessKey]struct{})
		stream := s.json.AcquireStream(nil)
		defer s.json.ReleaseStream(stream)

		if s.store != nil {
			for key, ptr := range s.store {
				typ, ok := s.getType(key.Action, key.Resource)

				if !ok {
					// This shouldn't happen
					continue
				}

				enc := s.json.EncoderOf(typ)
				stream.Reset(nil)
				enc.Encode(ptr, stream)

				if !yield(key, stream.Buffer()) {
					return
				}

				dumped[key.accessKey()] = struct{}{}
			}
		}

		for key, typ := range s.types {
			if _, skip := dumped[key]; skip {
				continue
			}

			ptr := reflect.New(typ).UnsafePointer()
			enc := s.json.EncoderOf(typ)
			stream.Reset(nil)
			enc.Encode(ptr, stream)

			if !yield(PolicyKey{Action: key.action, Resource: key.resource}, stream.Buffer()) {
				return
			}
		}
	}
}
