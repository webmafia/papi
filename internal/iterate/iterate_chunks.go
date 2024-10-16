package iterate

import (
	"iter"
	"strings"
)

func IterateChunks(s string, sep byte) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		var idx int

		for len(s) > 0 {
			i := strings.IndexByte(s, sep)

			if i < 0 {
				break
			}

			if len(s[:i]) > 0 {
				if !yield(idx, s[:i]) {
					return
				}

				idx++
			}

			s = s[i+1:]
		}

		if len(s) > 0 {
			yield(idx, s)
		}
	}
}
