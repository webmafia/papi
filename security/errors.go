package security

import "github.com/webmafia/papi/errors"

var (
	ErrInvalidAuthToken = errors.NewError("INVALID_TOKEN", "Invalid authentication token", 401)
	ErrAccessDenied     = errors.NewFrozenError("ACCESS_DENIED", "Access denied", 403)
)
