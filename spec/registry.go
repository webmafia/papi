package spec

import (
	"fmt"
	"reflect"

	"github.com/webmafia/fastapi/scanner"
	"github.com/webmafia/fastapi/scanner/structs"
	"github.com/webmafia/fastapi/scanner/value"
	"github.com/webmafia/fastapi/spec/openapi"
)

type Registry struct {
	schemas    map[reflect.Type]*openapi.Schema
	operations map[string][]*openapi.Operation
	opTagScan  value.ValueScanner
}

func NewRegistry(scanReg *scanner.Registry) (r *Registry, err error) {
	scan, err := structs.CreateTagScanner(scanReg, reflect.TypeOf((*inputTags)(nil)).Elem())

	if err != nil {
		return
	}

	r = &Registry{
		schemas:    make(map[reflect.Type]*openapi.Schema),
		operations: make(map[string][]*openapi.Operation),
		opTagScan:  scan,
	}

	return
}

func (r *Registry) RegisterSchema(typ reflect.Type, schema *openapi.Schema) {
	if schema == nil {
		delete(r.schemas, typ)
	} else {
		r.schemas[typ] = schema
	}
}

func (r *Registry) AddOperation(path string, op *openapi.Operation) (err error) {
	ops := r.operations[path]

	for _, o := range ops {
		if o.Id == op.Id {
			return fmt.Errorf("duplicate operation ID: %s", op.Id)
		}

		if o.Method == op.Method {
			return fmt.Errorf("duplicate method '%s' for path: %s", op.Method, path)
		}
	}

	r.operations[path] = append(ops, op)
	return
}

func (r *Registry) DescribeOperation(op *openapi.Operation, in, out reflect.Type) (err error) {

}
