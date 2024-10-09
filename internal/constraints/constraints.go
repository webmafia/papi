package constraints

// Signed is a constraint that permits any signed integer type.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint that permits any unsigned integer type.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Integer is a constraint that permits any integer type.
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint that permits any floating-point type.
type Float interface {
	~float32 | ~float64
}

// Number is a constraint that permits any numeric type.
type Number interface {
	Integer | Float
}

// Ordered is a constraint that permits any ordered type: any type
// that supports the operators < <= >= >.
type Ordered interface {
	Number | ~string
}
