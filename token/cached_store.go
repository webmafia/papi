package token

import (
	"context"

	"github.com/webmafia/identifier"
	"github.com/webmafia/lru"
)

var _ TokenStore = (*cachedStore)(nil)

type cachedStore struct {
	cache lru.LRU[identifier.ID, User]
	store TokenStore
}

func NewCachedStore(store TokenStore, capacity int) TokenStore {

	// Abort if already cached
	if s, ok := store.(*cachedStore); ok {
		return s
	}

	return &cachedStore{
		cache: lru.NewThreadSafe[identifier.ID, User](capacity),
		store: store,
	}
}

// Lookup in cache, and then in store.
func (c *cachedStore) Lookup(ctx context.Context, tok Token) (user User, err error) {
	return c.cache.GetOrSet(tok.Id(), func(_ identifier.ID) (User, error) {
		return c.store.Lookup(ctx, tok)
	})
}

// Store in cache and store.
func (c *cachedStore) Store(ctx context.Context, tok Token) error {
	c.cache.Remove(tok.Id())
	return c.store.Store(ctx, tok)
}
