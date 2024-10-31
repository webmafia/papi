package security

import (
	"context"

	"github.com/webmafia/identifier"
	"github.com/webmafia/lru"
)

type CachedTokenStore interface {
	TokenStore
	ClearCache()
}

var _ TokenStore = (*cachedStore)(nil)

type cachedStore struct {
	cache lru.LRU[identifier.ID, User]
	store TokenStore
}

// A cached store is used to reduce the preasure on the underlying store, and decrease
// any latency.
func NewCachedStore(store TokenStore, capacity int) CachedTokenStore {

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
func (c *cachedStore) Insert(ctx context.Context, tok Token) error {
	c.cache.Remove(tok.Id())
	return c.store.Insert(ctx, tok)
}

func (c *cachedStore) Delete(ctx context.Context, tokId identifier.ID) error {
	c.cache.Remove(tokId)
	return c.store.Delete(ctx, tokId)
}

// Clear token cache (without affecting underlying store)
func (c *cachedStore) ClearCache() {
	c.cache.Reset()
}
