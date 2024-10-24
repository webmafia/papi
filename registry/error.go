package registry

import "github.com/webmafia/papi/errors"

var (
	ErrFailedDecodeBody  = errors.NewFrozenError("FAILED_DECODE_BODY", "Could not decode JSON body", 400)
	ErrFailedDecodeParam = errors.NewFrozenError("FAILED_DECODE_PARAM", "Could not decode URL param", 400)
	ErrFailedDecodeQuery = errors.NewFrozenError("FAILED_DECODE_QUERY", "Could not decode query param", 400)
)
