package token

import (
	"context"
)

var (
	_ TokenStore = dummyStore{}
	_ User       = dummyUser{}
)

func DummyStore(roles ...string) TokenStore {
	return dummyStore{
		user: dummyUser{
			roles: roles,
		},
	}
}

type dummyStore struct {
	user dummyUser
}

func (d dummyStore) Lookup(ctx context.Context, tok Token) (user User, err error) {
	return d.user, nil
}

func (dummyStore) Store(ctx context.Context, tok Token) error {
	return nil
}

type dummyUser struct {
	roles []string
}

func (d dummyUser) UserRoles() []string {
	return d.roles
}
