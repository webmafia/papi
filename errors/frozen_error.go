package errors

import (
	jsoniter "github.com/json-iterator/go"
)

// An immutable error used for spawning new "explained" errors.
type FrozenError interface {
	ErrorDocumentor
	Explained(location, expect string) Error
}

type frozenError struct {
	status  int    // HTTP status code (e.g. 400)
	code    string // Error code (e.g. "TOO_SHORT")
	message string // Error message (e.g. "Too short")
}

// Create an immutable ("frozen") error that is used for spawning new "explained" errors.
func NewFrozenError(code string, message string, statusCode ...int) FrozenError {
	c := 400

	if len(statusCode) > 0 {
		c = statusCode[0]
	}

	return &frozenError{
		status:  c,
		code:    code,
		message: message,
	}
}

// Returns an `Error` with additional information.
func (err *frozenError) Explained(location, expect string) Error {
	return Error{
		status:   err.status,
		code:     err.code,
		message:  err.message,
		location: location,
		expect:   expect,
	}
}

func (err *frozenError) Status() int {
	return err.status
}

func (err *frozenError) JsonEncode(s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("code")
	s.WriteString(err.code)

	s.WriteMore()
	s.WriteObjectField("message")
	s.WriteString(err.message)

	s.WriteObjectEnd()
}

func (err *frozenError) ErrorDocument(s *jsoniter.Stream) {
	s.WriteObjectStart()
	s.WriteObjectField("errors")
	s.WriteArrayStart()
	err.JsonEncode(s)
	s.WriteArrayEnd()
	s.WriteObjectEnd()
}
