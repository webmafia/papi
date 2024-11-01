package security

import (
	"bytes"
	"errors"
	"fmt"
	"iter"
	"math"
	"reflect"
	"sync"
	"unsafe"

	"github.com/webmafia/papi/internal/json"
)

type PolicyKey struct {
	Role string
	Perm Permission
}

type Policy struct {
	Prio int64
	Cond []byte
}

type policyStore struct {
	store map[PolicyKey]policy
	types map[Permission]reflect.Type
	mu    sync.RWMutex
}

type policy struct {
	prio int64
	cond unsafe.Pointer
}

func (s *policyStore) Register(perm Permission, typ reflect.Type) (err error) {
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

	if s.types == nil {
		s.types = make(map[Permission]reflect.Type)
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

func (s *policyStore) freeze() {
	if s.store == nil {
		s.store = make(map[PolicyKey]policy)
	}
}

func (s *policyStore) Add(role string, perm Permission, prio int64, cond ...any) (err error) {
	typ, ok := s._GetType(perm)

	if !ok {
		return fmt.Errorf("no policy for %s found", perm)
	}

	var ptr unsafe.Pointer

	if len(cond) > 0 {
		switch c := cond[0].(type) {

		case []byte:
			var read bytes.Reader
			read.Reset(c)

			dec := json.DecoderOf(typ)
			ptr = reflect.New(typ).UnsafePointer()
			iter := json.AcquireIterator(&read)
			defer json.ReleaseIterator(iter)
			dec.Decode(ptr, iter)

			if iter.Error != nil {
				return iter.Error
			}

		default:
			cVal := reflect.ValueOf(cond)
			cTyp := cVal.Type()

			if cTyp.Kind() != reflect.Pointer {
				return errors.New("policy must be eiher a byte slice or a pointer")
			}

			if cTyp != typ && cTyp.Elem() != typ {
				return fmt.Errorf("invalid policy type: expected %s, but got %s", typ, cTyp)
			}

			ptr = cVal.UnsafePointer()
		}
	}

	s._Set(role, perm, prio, ptr)
	return
}

func (s *policyStore) BatchAdd(cb func(add func(role string, perm Permission, prio int64, condJson []byte) error) error) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.freeze()

	var read bytes.Reader
	iter := json.AcquireIterator(nil)
	defer json.ReleaseIterator(iter)

	add := func(role string, perm Permission, prio int64, condJson []byte) error {
		read.Reset(condJson)
		iter.Reset(&read)

		typ, ok := s.getType(perm)

		if !ok {
			return fmt.Errorf("no policy for %s found", perm)
		}

		dec := json.DecoderOf(typ)
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

func (s *policyStore) _GetType(perm Permission) (typ reflect.Type, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.getType(perm)
}

func (s *policyStore) getType(perm Permission) (typ reflect.Type, ok bool) {
	typ, ok = s.types[perm]
	return
}

func (s *policyStore) _Set(role string, perm Permission, prio int64, cond unsafe.Pointer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.freeze()
	s.set(role, perm, prio, cond)
}

func (s *policyStore) set(role string, perm Permission, prio int64, cond unsafe.Pointer) {
	s.store[PolicyKey{Role: role, Perm: perm}] = policy{
		prio: prio,
		cond: cond,
	}
}

func (s *policyStore) IteratePolicies() iter.Seq2[PolicyKey, Policy] {
	return func(yield func(PolicyKey, Policy) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		if s.store == nil {
			return
		}

		stream := json.AcquireStream(nil)
		defer json.ReleaseStream(stream)

		for key, pol := range s.store {
			typ, ok := s.getType(key.Perm)

			if !ok {
				// This shouldn't happen
				continue
			}

			enc := json.EncoderOf(typ)
			stream.Reset(nil)
			enc.Encode(pol.cond, stream)

			if !yield(key, Policy{Prio: pol.prio, Cond: stream.Buffer()}) {
				return
			}
		}
	}
}

func (s *policyStore) IteratePermissions(inPolicy ...bool) iter.Seq[Permission] {
	return func(yield func(Permission) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		var used, unused bool

		if len(inPolicy) > 0 {
			if inPolicy[0] {
				used = true
			} else {
				unused = true
			}
		} else {
			used = true
			unused = true
		}

		did := make(map[Permission]struct{})

		if s.store != nil {
			for key := range s.store {
				if used {
					if _, ok := did[key.Perm]; !ok && !yield(key.Perm) {
						return
					}
				}

				did[key.Perm] = struct{}{}
			}
		}

		if unused {
			for perm := range s.types {
				if _, ok := did[perm]; ok {
					continue
				}

				if !yield(perm) {
					return
				}
			}
		}
	}
}

func (s *policyStore) Get(roles []string, perm Permission) (cond unsafe.Pointer, err error) {
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

func (s *policyStore) get(role string, perm Permission) (pol policy, found bool) {
	pol, found = s.store[PolicyKey{Role: role, Perm: perm}]
	return
}

func (s *policyStore) Remove(role string, perm Permission) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.remove(role, perm)
}

func (s *policyStore) remove(role string, perm Permission) {
	if s.store != nil {
		delete(s.store, PolicyKey{Role: role, Perm: perm})
	}
}
