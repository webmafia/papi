package iterate

import (
	"iter"
)

func IterateStructTags[T ~string](tag T) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for tag != "" {
			// Skip leading space.
			i := 0
			for i < len(tag) && tag[i] == ' ' {
				i++
			}
			tag = tag[i:]
			if tag == "" {
				break
			}

			// Scan to colon. A space, a quote or a control character is a syntax error.
			// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
			// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
			// as it is simpler to inspect the tag's bytes than the tag's runes.
			i = 0
			for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
				i++
			}
			if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
				break
			}
			name := string(tag[:i])
			tag = tag[i+1:]

			// Scan quoted string to find value.
			i = 1
			for i < len(tag) && tag[i] != '"' {
				if tag[i] == '\\' {
					i++
				}
				i++
			}
			if i >= len(tag) {
				break
			}
			qvalue := string(tag[:i+1])
			tag = tag[i+1:]

			if l := len(qvalue); l >= 2 {
				if !yield(name, qvalue[1:l-1]) {
					break
				}
			}
		}
	}
}
