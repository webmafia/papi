package errors

import (
	"fmt"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type Errors []Error

var (
	_ error           = Errors{}
	_ ErrorDocumentor = Errors{}
	_ fmt.Stringer    = Errors{}
)

func (errs *Errors) Append(err Error) {
	*errs = append(*errs, err)
}

func (errs *Errors) Merge(errors Errors) {
	*errs = append(*errs, errors...)
}

func (errs *Errors) Reset() {
	for i := range *errs {
		(*errs)[i].Reset()
	}

	*errs = (*errs)[:0]
}

func (errs Errors) Len() int {
	return len(errs)
}

func (errs Errors) HasError() bool {
	return len(errs) > 0
}

func (errs Errors) Error() string {
	return errs.String()
}

func (errs Errors) String() string {
	if !errs.HasError() {
		return "(no error)"
	}

	var b strings.Builder

	for i := range errs {
		if i != 0 {
			b.WriteString("\n")
		}

		if errs[i].location != "" {
			b.WriteString(errs[i].location)
			b.WriteString(" - ")
		}

		b.WriteString(errs[i].code)
		b.WriteString(": ")
		b.WriteString(errs[i].message)

		if errs[i].expect != "" {
			b.WriteString(" (")
			b.WriteString(errs[i].expect)
			b.WriteString(")")
		}
	}

	return b.String()
}

func (errs Errors) Status() int {
	if !errs.HasError() {
		return 200
	}

	return errs[0].status
}

func (errs Errors) JsonEncode(s *jsoniter.Stream) {
	s.WriteArrayStart()

	for i := range errs {
		if i != 0 {
			s.WriteMore()
		}

		errs[i].JsonEncode(s)
	}

	s.WriteArrayEnd()
}

func (errs Errors) ErrorDocument(s *jsoniter.Stream) {
	s.WriteObjectStart()
	s.WriteObjectField("errors")
	errs.JsonEncode(s)
	s.WriteObjectEnd()
}
