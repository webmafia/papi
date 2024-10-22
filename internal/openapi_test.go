package internal

import (
	"testing"

	"github.com/webmafia/papi/openapi"
)

// Ensure that the local `document` is exactly the same as `openapi.Document`.
func Test_document(t *testing.T) {
	if !EqualStructs[document, openapi.Document]() {
		t.Errorf("%T and %T mismatch", document{}, openapi.Document{})
	}
}
