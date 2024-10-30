package policy

import (
	"bytes"
	"errors"
	"fmt"
	"iter"
	"math"
	"reflect"
	"sync"
	"unsafe"

	"github.com/webmafia/papi/internal"
)

type Store struct {
	store map[PolicyKey]policy
	types map[Permission]reflect.Type
	json  *internal.JSONPool
	mu    sync.RWMutex
}

type policy struct {
	prio int64
	cond unsafe.Pointer
}

type Policy struct {
	Prio int64
	Cond []byte
}

func NewStore(json *internal.JSONPool) *Store {
	return &Store{
		types: make(map[Permission]reflect.Type),
		json:  json,
	}
}

func (s *Store) Register(perm Permission, typ reflect.Type) (err error) {
	if !perm.HasAction() {
		return errors.New("missing 'action' for policy")
	}

	if !perm.HasResource() {
		return errors.New("missing 'resource' for policy")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.store != nil {
		return errors.New("store has already been frozen and can't accept any new registration")
	}

	if regTyp, ok := s.types[perm]; ok {
		if regTyp != typ {
			err = fmt.Errorf("can't register %s as policy for %s, as %s already is registered", typ, perm, regTyp)
		}

		return
	}

	s.types[perm] = typ
	return
}

func (s *Store) Freeze() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.freeze()
}

func (s *Store) freeze() {
	if s.store == nil {
		s.store = make(map[PolicyKey]policy)
	}
}

func (s *Store) Add(role string, perm Permission, prio int64, condJson []byte) (err error) {
	typ, ok := s._GetType(perm)

	if !ok {
		return fmt.Errorf("no policy for %s found", perm)
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

	s._Set(role, perm, prio, ptr)
	return
}

func (s *Store) BatchAdd(cb func(add func(role string, perm Permission, prio int64, condJson []byte) error) error) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.freeze()

	var read bytes.Reader
	iter := s.json.AcquireIterator(nil)
	defer s.json.ReleaseIterator(iter)

	add := func(role string, perm Permission, prio int64, condJson []byte) error {
		read.Reset(condJson)
		iter.Reset(&read)

		typ, ok := s.getType(perm)

		if !ok {
			return fmt.Errorf("no policy for %s found", perm)
		}

		dec := s.json.DecoderOf(typ)
		ptr := reflect.New(typ).UnsafePointer()

		dec.Decode(ptr, iter)

		if iter.Error != nil {
			return iter.Error
		}

		s.set(role, perm, prio, ptr)
		return nil
	}

	return cb(add)
}

func (s *Store) _GetType(perm Permission) (typ reflect.Type, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.getType(perm)
}

func (s *Store) getType(perm Permission) (typ reflect.Type, ok bool) {
	typ, ok = s.types[perm]
	return
}

func (s *Store) _Set(role string, perm Permission, prio int64, cond unsafe.Pointer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.freeze()
	s.set(role, perm, prio, cond)
}

func (s *Store) set(role string, perm Permission, prio int64, cond unsafe.Pointer) {
	s.store[PolicyKey{Role: role, Perm: perm}] = policy{
		prio: prio,
		cond: cond,
	}
}

func (s *Store) Iterate() iter.Seq2[PolicyKey, Policy] {
	return func(yield func(PolicyKey, Policy) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		dumped := make(map[Permission]struct{})
		stream := s.json.AcquireStream(nil)
		defer s.json.ReleaseStream(stream)

		if s.store != nil {
			for key, pol := range s.store {
				typ, ok := s.getType(key.Perm)

				if !ok {
					// This shouldn't happen
					continue
				}

				enc := s.json.EncoderOf(typ)
				stream.Reset(nil)
				enc.Encode(pol.cond, stream)

				if !yield(key, Policy{Prio: pol.prio, Cond: stream.Buffer()}) {
					return
				}

				dumped[key.Perm] = struct{}{}
			}
		}

		for perm, typ := range s.types {
			if _, skip := dumped[perm]; skip {
				continue
			}

			ptr := reflect.New(typ).UnsafePointer()
			enc := s.json.EncoderOf(typ)
			stream.Reset(nil)
			enc.Encode(ptr, stream)

			if !yield(PolicyKey{Perm: perm}, Policy{Cond: stream.Buffer()}) {
				return
			}
		}
	}
}

func (s *Store) Get(roles []string, perm Permission) (cond unsafe.Pointer, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var prio int64 = math.MaxInt64
	var found bool

	for _, role := range roles {
		pol, ok := s.get(role, perm)

		if ok && pol.prio < prio {
			found = true
			prio = pol.prio
			cond = pol.cond
		}
	}

	if !found {
		return nil, ErrAccessDenied.Detailed(fmt.Sprintf("Missing permission %s", perm))
	}

	return
}

func (s *Store) get(role string, perm Permission) (pol policy, found bool) {
	pol, found = s.store[PolicyKey{Role: role, Perm: perm}]
	return
}
