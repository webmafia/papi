package valid

import (
	"testing"

	"github.com/webmafia/papi/errors"
)

func BenchmarkValidateStruct(b *testing.B) {
	var v struct{}
	var errs errors.Errors

	for b.Loop() {
		_ = ValidateStruct(&v, &errs)
	}
}
