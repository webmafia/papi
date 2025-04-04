package token

import (
	"context"
)

var _ TokenStore = dummyStore{}

// Used for testing.
func DummyStore(roles ...string) TokenStore {
	return dummyStore{
		roles: roles,
	}
}

type dummyStore struct {
	roles []string
}

func (d dummyStore) UserRoles(_ context.Context, _ Token) (roles []string, err error) {
	return d.roles, nil
}
