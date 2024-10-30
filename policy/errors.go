package policy

import (
	"github.com/webmafia/papi/errors"
)

var ErrAccessDenied = errors.NewFrozenError("ACCESS_DENIED", "Access denied", 403)
