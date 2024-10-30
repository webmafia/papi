package token

import "context"

type TokenStore interface {
	Lookup(ctx context.Context, tok Token) (user User, err error)
	Store(ctx context.Context, tok Token) error
}

type User interface {
	UserRoles() []string
}
