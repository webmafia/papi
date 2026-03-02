package internal

import (
	"iter"
	"sync"
)

// A thread-safe set.
type Set[T comparable] struct {
	perms map[T]struct{}
	mu    sync.Mutex
}

func (s *Set[T]) Add(v T) (existed bool) {
	s.mu.Lock()

	if s.perms == nil {
		s.perms = make(map[T]struct{})
	} else {
		_, existed = s.perms[v]
	}

	s.perms[v] = struct{}{}
	s.mu.Unlock()

	return
}

func (s *Set[T]) Clear() {
	s.mu.Lock()

	if s.perms != nil {
		clear(s.perms)
	}

	s.mu.Unlock()
}

func (s *Set[T]) Delete(v T) (existed bool) {
	s.mu.Lock()

	if s.perms != nil {
		_, existed = s.perms[v]
		delete(s.perms, v)
	}

	s.mu.Unlock()

	return
}

func (s *Set[T]) Has(v T) (exists bool) {
	s.mu.Lock()

	if s.perms != nil {
		_, exists = s.perms[v]
	}

	s.mu.Unlock()

	return
}

func (s *Set[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		s.mu.Lock()
		defer s.mu.Unlock()

		if s.perms == nil {
			return
		}

		for v := range s.perms {
			if !yield(v) {
				return
			}
		}
	}
}
