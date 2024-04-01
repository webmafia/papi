package fastapi

// func createHandler[U, I, O any](r Route[U, I, O]) (h fasthttp.RequestHandler, err error) {
// 	var (
// 		in  I
// 		out O
// 	)

// 	typ := reflect.TypeOf(in)

// 	if typ.Kind() != reflect.Struct {
// 		return nil, fmt.Errorf("invalid input type; expected struct, got %s", typ.Kind())
// 	}

// 	num := typ.NumField()

// 	for i := 0; i < num; i++ {
// 		f := typ.Field(i)

// 		if !f.IsExported() {
// 			continue
// 		}
// 	}

// 	return func(ctx *fasthttp.RequestCtx) {
// 		ctx.UserValue()
// 	}, nil
// }
