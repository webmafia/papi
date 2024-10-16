package papi

import "github.com/webbmaffian/papi/errors"

var (
	ErrNotFound      = errors.NewFrozenError("NOT_FOUND", "The API route could not be found", 404)
	ErrInvalidParams = errors.NewFrozenError("INVALID_PARAMS", "URL params count mismatch", 500)
	ErrUnknownError  = errors.NewFrozenError("UNKNOWN_ERROR", "An unknown error occured", 500)
)
