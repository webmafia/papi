package token

import (
	"context"
)

var (
	_ TokenStore = dummyStore{}
	_ User       = dummyUser{}
)

// Used for testing.
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

func (d dummyStore) Lookup(_ context.Context, _ Token) (user User, err error) {
	return d.user, nil
}

type dummyUser struct {
	roles []string
}

func (d dummyUser) UserRoles() []string {
	return d.roles
}
