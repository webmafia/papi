package internal

// Panics if err != nil.
//
//go:inline
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}
