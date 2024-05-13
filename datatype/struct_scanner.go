package datatype

func CreateStructScanner[T any](d *DataTypes, tags ...string) (fn func(*T, func(tag, val string) string), err error) {
	// TODO: Build slice of props (tag key, tag value, scanner) and return
	// a function that iterates it and calls callback for each key-value pair,
	// scanning its returning string.

	return
}
