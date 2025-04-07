package token

import (
	"context"
	"time"
)

type Store[T any] interface {

	// Looks up a token in the underlying token store, and returns its corresponding user roles.
	// A user can have 0+ roles. If the token doesn't exist in store and/or has been revoked, it MUST
	// return an error. The ctx MIGHT be a *papi.RequestCtx.
	UserRoles(ctx context.Context, tok Token) (roles []string, err error)

	// Consume an authentication code and returns its corresponding details.
	ConsumeAuthCode(ctx context.Context, code OneTimeCode) (userId T, cookie bool, err error)

	// Save an authentication code to a storage.
	SaveAuthCode(ctx context.Context, userId T, code OneTimeCode, expiry time.Time, cookie bool) error

	// Save an access token to a storage. Only the ID of the token should be saved, not the whole token.
	SaveAccessToken(ctx context.Context, userId T, tok Token, cookie bool) error
}
