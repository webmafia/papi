package papi

import (
	errs "errors"

	"github.com/webbmaffian/papi/errors"
)

// API error
var (
	ErrNotFound      = errors.NewFrozenError("NOT_FOUND", "The API route could not be found", 404)
	ErrInvalidParams = errors.NewFrozenError("INVALID_PARAMS", "URL params count mismatch", 500)
	ErrUnknownError  = errors.NewFrozenError("UNKNOWN_ERROR", "An unknown error occured", 500)
)

var (
	ErrInvalidOpenAPI      = errs.New("there must not be any existing operations in OpenAPI documentation")
	ErrMissingOpenAPI      = errs.New("no OpenAPI documentation initialized")
	ErrMissingRoutePath    = errs.New("missing route path")
	ErrMissingRouteHandler = errs.New("missing route handler")
	ErrMissingRouteMethod  = errs.New("missing route method")
)
