package token

import (
	"context"
	"time"

	"github.com/webmafia/papi/internal"
	"github.com/webmafia/papi/security"
)

type Store[T any] interface {

	// Checks whether a token has a permission, and sets any conditions. Any returned error will result in "403 Forbidden".
	CheckPermission(ctx context.Context, tokId uint64, perm security.Permission, cond internal.Setter) (err error)

	// Looks up a token in the underlying token store, and returns its corresponding user roles.
	// A user can have 0+ roles. If the token doesn't exist in store and/or has been revoked, it MUST
	// return an error. The ctx MIGHT be a *papi.RequestCtx.
	// UserRoles(ctx context.Context, tokId uint64) (roles []string, err error)

	// Consume an authentication code and returns its corresponding details.
	ConsumeAuthCode(ctx context.Context, code string) (userId T, cookie bool, err error)

	// Save an authentication code to a storage.
	SaveAuthCode(ctx context.Context, userId T, code string, expiry time.Time, cookie bool) error

	// Save an access token to a storage. Only the ID of the token should be saved, not the whole token.
	SaveAccessToken(ctx context.Context, userId T, tokId uint64, cookie bool) error
}
