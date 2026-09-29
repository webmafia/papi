package registry

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/valyala/fasthttp"
	"github.com/webmafia/papi/openapi"
)

type queryIgnoreFilters struct {
	Included bool `query:"included"`
	Ignored1 bool `query:"-"`
	Ignored2 bool `query:"-"`
}

type queryIgnoreRequest struct {
	Filters queryIgnoreFilters
}

func TestDescribeOperationIgnoresDashQueryTags(t *testing.T) {
	reg := NewRegistry()
	op := &openapi.Operation{}

	if err := reg.DescribeOperation(
		op,
		reflect.TypeFor[queryIgnoreRequest](),
		reflect.TypeFor[int](),
	); err != nil {
		t.Fatal(err)
	}

	if len(op.Parameters) != 1 {
		t.Fatalf("expected one query parameter, got %d", len(op.Parameters))
	}

	if name := op.Parameters[0].Name; name != "included" {
		t.Fatalf("expected included query parameter, got %q", name)
	}
}

func TestCreateBinderIgnoresDashQueryTags(t *testing.T) {
	reg := NewRegistry()
	bind, _, err := reg.CreateBinder(reflect.TypeFor[queryIgnoreRequest](), nil, nil)

	if err != nil {
		t.Fatal(err)
	}

	ctx := &fasthttp.RequestCtx{}
	ctx.QueryArgs().Set("included", "true")
	ctx.QueryArgs().Set("-", "true")

	var req queryIgnoreRequest

	if err = bind(ctx, unsafe.Pointer(&req)); err != nil {
		t.Fatal(err)
	}

	if !req.Filters.Included {
		t.Error("expected included query parameter to be bound")
	}

	if req.Filters.Ignored1 || req.Filters.Ignored2 {
		t.Error("expected dash query tags to be ignored")
	}
}
