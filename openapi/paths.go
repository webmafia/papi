package openapi

import "fmt"

type Paths map[string][]*Operation

func (p Paths) AddOperation(path string, op *Operation) (err error) {
	for opPath, ops := range p {
		for _, o := range ops {

			// Operation ID must be unique among all paths
			// if o.Id == op.Id {
			// 	return fmt.Errorf("duplicate operation ID: %s", op.Id)
			// }

			// Method must be unique for specific path
			if opPath == path && o.Method == op.Method {
				return fmt.Errorf("duplicate method '%s' for path: %s", op.Method, path)
			}
		}
	}

	p[path] = append(p[path], op)
	return
}

// func (p Paths) JsonEncode(ctx *encoderContext, s *jsoniter.Stream) {
// 	m := make(map[string][]Operation)

// 	for i := range p {
// 		m[p[i].Path] = append(m[p[i].Path], p[i])
// 	}

// 	s.WriteObjectStart()

// 	var written bool

// 	for path, ops := range m {
// 		if written {
// 			s.WriteMore()
// 		} else {
// 			written = true
// 		}

// 		s.WriteObjectField(path)
// 		s.WriteObjectStart()

// 		for i := range ops {
// 			if i != 0 {
// 				s.WriteMore()
// 			}

// 			s.WriteObjectField(strings.ToLower(ops[i].Method))
// 			ops[i].JsonEncode(ctx, s)
// 		}

// 		s.WriteObjectEnd()
// 	}

// 	s.WriteObjectEnd()
// }
