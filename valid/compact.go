package valid

// Ensures that the capacity of a slice is exactly the length
func compactSlice[T any](s []T) []T {
	if l := len(s); cap(s) > l {
		s2 := make([]T, l)
		copy(s2, s)
		return s2
	}

	return s
}
