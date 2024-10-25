package valid

import (
	"github.com/webmafia/papi/errors"
)

var (
	ErrTooLow        = errors.NewFrozenError("TOO_LOW", "Value is too low", 400)
	ErrTooHigh       = errors.NewFrozenError("TOO_HIGH", "Value is too high", 400)
	ErrTooShort      = errors.NewFrozenError("TOO_SHORT", "Value is too short", 400)
	ErrTooLong       = errors.NewFrozenError("TOO_LONG", "Value is too long", 400)
	ErrRequired      = errors.NewFrozenError("REQUIRED", "Value is required", 400)
	ErrFailedEnum    = errors.NewFrozenError("FAILED_ENUM", "Value is too low", 400)
	ErrFailedPattern = errors.NewFrozenError("FAILED_PATTERN", "Value does not match pattern", 400)
	ErrFailedDefault = errors.NewFrozenError("FAILED_DEFAULT", "Default value can't be set", 500)
)
