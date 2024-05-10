package fastapi

import (
	"github.com/gosimple/slug"
	"github.com/webmafia/fastapi/spec"
)

func addRouteDocs[U, I, O any](api *API[U], r Route[U, I, O]) (err error) {
	if api.docs == nil {
		return
	}

	op := spec.Operation{
		Id:          slug.Make(r.Summary),
		Path:        r.Path,
		Method:      string(r.Method),
		Summary:     r.Summary,
		Description: r.Description,
		Tags:        r.Tags,
	}

	var in I

	if err = op.ParamsFromStruct(in); err != nil {
		return
	}

	api.docs.Paths = append(api.docs.Paths, op)

	return
}
