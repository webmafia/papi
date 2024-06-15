package schema

import (
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

type Schema interface {
	Name() string
	ScanTags(tags reflect.StructTag) error
	EncodeSchema(*jsoniter.Stream)
}
