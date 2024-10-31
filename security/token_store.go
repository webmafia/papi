package security

import (
	"context"

	"github.com/webmafia/identifier"
)

type TokenStore interface {

	// Validates a token, looks it up in the underlying token store, and returns its corresponding user.
	// A user can have 0+ roles. If the token doesn't exist in store and/or has been revoked, it MUST
	// return an error. The ctx MIGHT be a *papi.RequestCtx.
	Lookup(ctx context.Context, tok Token) (user User, err error)

	// Inserts a newly created token into the store. This should NOT be called manually, as it's called
	// automatically after a token is created. The store should NOT save the whole token, only its ID, any
	// relation to its corresponding user, and any additional data that that might help later recovation.
	// The ctx MIGHT be a *papi.RequestCtx.
	Insert(ctx context.Context, tok Token) error

	// Deletes a token permanently. Any failure MUST return an error. A deleted token must NOT be able to
	// be looked up later. The ctx MIGHT be a *papi.RequestCtx.
	Delete(ctx context.Context, tokId identifier.ID) error
}

// A user representation. The representation SHOULD be brief and only contain values that are necessary to
// comply with policies, e.g. user ID, tenant ID, roles, etc.
type User interface {
	UserRoles() []string
}
