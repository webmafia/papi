package token

import (
	"context"
)

var _ TokenStore = dummyStore{}

type dummyStore struct{}

func (dummyStore) Lookup(ctx context.Context, tok Token) (user User, err error) {
	return
}

func (dummyStore) Store(ctx context.Context, tok Token) error {
	return nil
}
