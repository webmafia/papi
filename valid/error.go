package valid

import "fmt"

var (
	ErrBelowMin       ValidationError = &validationError{code: "BELOW_MIN", msg: "Below minimum"}
	ErrAboveMax       ValidationError = &validationError{code: "ABOVE_MAX", msg: "Above maximum"}
	ErrRequired       ValidationError = &validationError{code: "REQUIRED", msg: "Value is required"}
	ErrInvalidEnum    ValidationError = &validationError{code: "INVALID_ENUM", msg: "Value is invalid"}
	ErrInvalidPattern ValidationError = &validationError{code: "INVALID_PATTERN", msg: "Value is invalid"}
)

type ValidationError interface {
	error
	fmt.Stringer
	Code() string
}

type validationError struct {
	code string
	msg  string
}

func (v *validationError) Error() string {
	return v.msg
}

func (v *validationError) String() string {
	return v.msg
}

func (v *validationError) Code() string {
	return v.code
}

type FieldError struct {
	err    ValidationError
	field  string
	expect string
}

func (f FieldError) Error() string {
	return f.err.Error()
}

func (f FieldError) String() string {
	return f.err.String()
}

func (f FieldError) Code() string {
	return f.err.Code()
}

func (f FieldError) Field() string {
	return f.field
}

func (f FieldError) Expect() string {
	return f.expect
}

type FieldErrors []FieldError

func (f *FieldErrors) Append(err FieldError) {
	*f = append(*f, err)
}

func (f FieldErrors) HasError() bool {
	return len(f) > 0
}

func (f FieldErrors) Error() string {
	if f.HasError() {
		return f[0].Error()
	}

	return ""
}
