package token

import "github.com/webmafia/papi/errors"

var ErrInvalidAuthToken = errors.NewError("INVALID_TOKEN", "Invalid authentication token", 401)
