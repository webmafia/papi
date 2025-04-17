package security

import "github.com/webmafia/papi/errors"

var (
	ErrInvalidAuthToken = errors.NewError("INVALID_TOKEN", "Invalid authentication token", 401)
	ErrInvalidAuthCode  = errors.NewError("INVALID_CODE", "Invalid authentication code", 401)
	ErrAccessDenied     = errors.NewFrozenError("ACCESS_DENIED", "Access denied", 403)
)
