package fastapi

import (
	"github.com/webmafia/fastapi/internal"
	"github.com/webmafia/fastapi/scanner/strings"
	"github.com/webmafia/fastapi/scanner/structs"
)

func registerScanners(f *strings.Factory) (err error) {
	typ := internal.ReflectType[inputTags]()
	scan, err := structs.CreateTagScanner(f, typ)

	if err != nil {
		return
	}

	f.RegisterScanner(typ, scan)
	return
}
