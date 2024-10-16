package errors

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

var (
	_ error           = Error{}
	_ ErrorDocumentor = Error{}
	_ fmt.Stringer    = Error{}
)

type Error struct {
	status   int    // HTTP status code (e.g. 400)
	code     string // Error code (e.g. "TOO_SHORT")
	message  string // Error message (e.g. "Too short")
	location string // What the error concerns (e.g. the specific field "password")
	expect   string // What was expected (e.g. "7", as in "min 7 characters")
}

func NewError(code string, message string, statusCode ...int) Error {
	c := 400

	if len(statusCode) > 0 {
		c = statusCode[0]
	}

	return Error{
		status:  c,
		code:    code,
		message: message,
	}
}

func (err Error) JsonEncode(s *jsoniter.Stream) {
	s.WriteObjectStart()

	s.WriteObjectField("code")
	s.WriteString(err.code)

	s.WriteMore()
	s.WriteObjectField("message")
	s.WriteString(err.message)

	if err.location != "" {
		s.WriteMore()
		s.WriteObjectField("location")
		s.WriteString(err.location)
	}

	if err.expect != "" {
		s.WriteMore()
		s.WriteObjectField("expect")
		s.WriteString(err.expect)
	}

	s.WriteObjectEnd()
}

func (err Error) ErrorDocument(s *jsoniter.Stream) {
	s.WriteObjectStart()
	s.WriteObjectField("errors")
	s.WriteArrayStart()
	err.JsonEncode(s)
	s.WriteArrayEnd()
	s.WriteObjectEnd()
}

func (err Error) Error() string {
	return err.message
}

func (err Error) String() string {
	return err.message
}

func (err Error) Status() int {
	return err.status
}

func (err Error) Code() string {
	return err.code
}

func (err Error) Message() string {
	return err.message
}

func (err Error) Location() string {
	return err.location
}

func (err Error) Expect() string {
	return err.expect
}

func (err *Error) Reset() {
	err.status = 0
	err.code = ""
	err.message = ""
	err.location = ""
	err.expect = ""
}
