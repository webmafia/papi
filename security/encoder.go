package security

import "encoding/base32"

var encoder = base32.HexEncoding.WithPadding(base32.NoPadding)
