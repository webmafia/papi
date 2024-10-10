package valid

import "fmt"

var (
	_ fmt.Stringer = ValidationError("")
	_ error        = ValidationError("")
)

type ValidationError string

func (v ValidationError) String() string {
	return string(v)
}

func (v ValidationError) Error() string {
	return string(v)
}

func Error(s string, args ...any) ValidationError {
	if len(args) == 0 {
		return ValidationError(s)
	}

	return ValidationError(fmt.Sprintf(s, args...))
}
