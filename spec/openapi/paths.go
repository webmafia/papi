package openapi

type Paths map[string][]Operation

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
