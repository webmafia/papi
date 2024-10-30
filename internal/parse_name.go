package internal

import (
	"runtime"
	"strings"

	"github.com/webmafia/fast"
)

func CallerName(skip int) string {
	return CallerNameFromFunc(Caller(skip + 1))
}

func CallerNameFromFunc(f *runtime.Func) string {
	if f == nil {
		return ""
	}

	name := f.Name()

	i := strings.LastIndexByte(name, '.')

	if i > 0 {
		name = name[i+1:]
	}

	return name
}

func CallerType(skip int) string {
	return CallerTypeFromFunc(Caller(skip + 1))
}

func CallerTypeFromFunc(f *runtime.Func) string {
	if f == nil {
		return ""
	}

	name := f.Name()

	end := strings.LastIndexByte(name, '.')

	if end > 0 {
		name = name[:end]
		start := strings.LastIndexByte(name, '.')

		if start > 0 {
			name = name[start+1:]
		}
	}

	return name
}

func Caller(skip int) *runtime.Func {
	pc, _, _, ok := runtime.Caller(skip + 1)

	if !ok {
		return nil
	}

	return runtime.FuncForPC(pc)
}

func ParseName(s string) (title, operationId string) {
	b := fast.StringToBytes(s)
	l := len(b)
	alloc := calcAlloc(b)
	var tb strings.Builder
	var ob strings.Builder
	tb.Grow(alloc)
	ob.Grow(alloc)

	for i, c := range b {
		if i == 0 {
			tb.WriteByte(c)

			if isUpper(c) {
				ob.WriteByte(toLower(c))
			} else {
				ob.WriteByte(c)
			}
		} else if i != 0 && isUpper(c) && !isUpper(b[i-1]) {
			tb.WriteByte(' ')
			ob.WriteByte('-')
			ob.WriteByte(toLower(c))

			if i < l-1 && isUpper(b[i+1]) {
				tb.WriteByte(c)
			} else {
				tb.WriteByte(toLower(c))
			}
		} else if isAlphaNumeric(c) {
			tb.WriteByte(c)

			if isUpper(c) {
				ob.WriteByte(toLower(c))
			} else {
				ob.WriteByte(c)
			}
		} else {
			tb.WriteByte(' ')
			ob.WriteByte('-')
		}
	}

	return tb.String(), ob.String()
}

func calcAlloc(b []byte) (alloc int) {
	alloc += len(b)

	for i, c := range b {
		if i != 0 && isUpper(c) && !isUpper(b[i-1]) {
			alloc++
		}
	}

	return
}

//go:inline
func isUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

//go:inline
func isAlphaNumeric(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9')
}

func toLower(c byte) byte {
	return c - 'A' + 'a'
}
