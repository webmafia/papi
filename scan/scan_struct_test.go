package scan

import (
	"log"
	"reflect"
	"testing"
)

func Example_iterateStructTags() {
	type Foobar struct {
		Foo int    `json:"wazzaaa" param:"waaaaa" query:"zzaaaaa"`
		Bar string `json:"a" param:"b b b" query:"c"`
	}

	typ := reflect.TypeOf(Foobar{})
	l := typ.NumField()

	for i := 0; i < l; i++ {
		fld := typ.Field(i)

		iterateStructTags(fld.Tag, func(key, val string) bool {
			log.Println(key, "is", val)
			return false
		})
	}

	// Output: Todo
}

func Benchmark_iterateStructTags(b *testing.B) {
	type Foobar struct {
		Foo int    `json:"wazzaaa" param:"waaaaa" query:"zzaaaaa"`
		Bar string `json:"a" param:"b b b" query:"c"`
	}

	typ := reflect.TypeOf(Foobar{})
	fld := typ.Field(0)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = iterateStructTags(fld.Tag, func(key, val string) bool {
			_ = key
			_ = val
			return false
		})
	}
}
