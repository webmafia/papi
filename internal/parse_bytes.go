package internal

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

func ParseBytes(s string) (int64, error) {
	x := strings.TrimSpace(strings.ToLower(s))
	x = strings.ReplaceAll(strings.ReplaceAll(x, "_", ""), ",", "")
	if x == "" {
		return 0, errors.New("empty size")
	}

	// split numeric prefix and unit suffix
	i := 0
	for i < len(x) && (x[i] == '+' || x[i] == '-' || x[i] == '.' || (x[i] >= '0' && x[i] <= '9')) {
		i++
	}
	num, unit := x[:i], strings.TrimSpace(x[i:])
	if num == "" || strings.HasPrefix(num, "-") {
		return 0, errors.New("invalid number")
	}

	// normalize unit aliases (binary, 1024^n)
	var mul int64 = 1
	switch unit {
	case "", "b":
		mul = 1
	case "k", "kb", "kib":
		mul = 1 << 10
	case "m", "mb", "mib":
		mul = 1 << 20
	case "g", "gb", "gib":
		mul = 1 << 30
	case "t", "tb", "tib":
		mul = 1 << 40
	case "p", "pb", "pib":
		mul = 1 << 50
	default:
		return 0, errors.New("unknown unit")
	}

	// integer fast-path
	if !strings.ContainsAny(num, ".e") {
		u, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			return 0, err
		}
		if mul != 0 && u > math.MaxInt64/mul {
			return 0, errors.New("overflow")
		}
		return u * mul, nil
	}

	// decimal path
	f, err := strconv.ParseFloat(num, 64)
	if err != nil || f < 0 {
		return 0, errors.New("invalid number")
	}
	v := f * float64(mul)
	if v < 0 || v > float64(^uint64(0)) {
		return 0, errors.New("overflow")
	}
	return int64(v + 0.5), nil // round to nearest
}
