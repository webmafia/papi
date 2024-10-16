package iterate

import (
	"iter"
	"strings"
)

func IterateFlags(flags string) iter.Seq[string] {
	return func(yield func(string) bool) {
		str := flags

		for {
			i := strings.IndexByte(str, ',')

			if i < 0 {
				break
			}

			if len(str[:i]) > 0 {
				if !yield(str[:i]) {
					return
				}
			}

			str = str[i+1:]
		}

		if len(flags) > 0 {
			yield(str)
		}
	}
}

func HasFlag(flags string, flag string) bool {
	for f := range IterateFlags(flags) {
		if f == flag {
			return true
		}
	}

	return false
}
