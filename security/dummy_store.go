package security

import (
	"context"

	"github.com/webmafia/identifier"
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

func (d dummyStore) Lookup(ctx context.Context, tok Token) (user User, err error) {
	return d.user, nil
}

func (dummyStore) Insert(ctx context.Context, tok Token) error {
	return nil
}

func (d dummyStore) Delete(ctx context.Context, tokId identifier.ID) error {
	return nil
}

type dummyUser struct {
	roles []string
}

func (d dummyUser) UserRoles() []string {
	return d.roles
}
