package datatype

import (
	"fmt"
	"strconv"
	"strings"
)

func RegisterStandardScanners(d *DataTypes) {
	RegisterScanner(d, func(ptr *bool, str string) (err error) {
		switch str {
		case "1", "t", "T", "true", "TRUE", "True", "yes", "YES", "Yes":
			*ptr = true
		case "0", "f", "F", "false", "FALSE", "False", "no", "NO", "No":
			*ptr = false
		default:
			return fmt.Errorf("invalid boolean: '%s'", str)
		}

		return
	})

	RegisterScanner(d, func(ptr *complex64, str string) (err error) {
		v, err := strconv.ParseComplex(str, 64)

		if err == nil {
			*ptr = complex64(v)
		}

		return
	})

	RegisterScanner(d, func(ptr *complex128, str string) (err error) {
		v, err := strconv.ParseComplex(str, 128)

		if err == nil {
			*ptr = v
		}

		return
	})

	RegisterScanner(d, func(ptr *float32, str string) (err error) {
		v, err := strconv.ParseFloat(str, 32)
		if err == nil {
			*ptr = float32(v)
		}
		return
	})

	RegisterScanner(d, func(ptr *float64, str string) (err error) {
		v, err := strconv.ParseFloat(str, 64)
		if err == nil {
			*ptr = v
		}
		return
	})

	RegisterScanner(d, func(ptr *int, str string) (err error) {
		v, err := strconv.ParseInt(str, 10, 0)
		if err == nil {
			*ptr = int(v)
		}
		return
	})

	RegisterScanner(d, func(ptr *int8, str string) (err error) {
		v, err := strconv.ParseInt(str, 10, 8)
		if err == nil {
			*ptr = int8(v)
		}
		return
	})

	RegisterScanner(d, func(ptr *int16, str string) (err error) {
		v, err := strconv.ParseInt(str, 10, 16)
		if err == nil {
			*ptr = int16(v)
		}
		return
	})

	RegisterScanner(d, func(ptr *int32, str string) (err error) {
		v, err := strconv.ParseInt(str, 10, 32)
		if err == nil {
			*ptr = int32(v)
		}
		return
	})

	RegisterScanner(d, func(ptr *int64, str string) (err error) {
		v, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			*ptr = v
		}
		return
	})

	RegisterScanner(d, func(ptr *uint, str string) (err error) {
		v, err := strconv.ParseUint(str, 10, 0)
		if err == nil {
			*ptr = uint(v)
		}
		return
	})

	RegisterScanner(d, func(ptr *uint8, str string) (err error) {
		v, err := strconv.ParseUint(str, 10, 8)
		if err == nil {
			*ptr = uint8(v)
		}
		return
	})

	RegisterScanner(d, func(ptr *uint16, str string) (err error) {
		v, err := strconv.ParseUint(str, 10, 16)
		if err == nil {
			*ptr = uint16(v)
		}
		return
	})

	RegisterScanner(d, func(ptr *uint32, str string) (err error) {
		v, err := strconv.ParseUint(str, 10, 32)
		if err == nil {
			*ptr = uint32(v)
		}
		return
	})

	RegisterScanner(d, func(ptr *uint64, str string) (err error) {
		v, err := strconv.ParseUint(str, 10, 64)
		if err == nil {
			*ptr = v
		}
		return
	})

	RegisterScanner(d, func(ptr *string, str string) (err error) {
		*ptr = str
		return
	})

	RegisterScanner(d, func(ptr *string, str string) (err error) {
		*ptr = strings.Clone(str)
		return
	})
}
