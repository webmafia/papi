package security

import (
	"context"
)

type TokenStore interface {

	// Looks up a token in the underlying token store, and returns its corresponding user.
	// A user can have 0+ roles. If the token doesn't exist in store and/or has been revoked, it MUST
	// return an error. The ctx MIGHT be a *papi.RequestCtx.
	Lookup(ctx context.Context, tok Token) (user User, err error)
}
