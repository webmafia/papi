package internal

import (
	"unsafe"

	"github.com/webbmaffian/papi/openapi"
)

// Must be synchronized with `openapi.Document`.
type document struct {
	_     openapi.Info
	_     []openapi.Server
	paths openapi.Paths
}

func AddOperationToDocument(doc *openapi.Document, path string, op *openapi.Operation) (err error) {
	return (*document)(unsafe.Pointer(doc)).paths.AddOperation(path, op)
}
