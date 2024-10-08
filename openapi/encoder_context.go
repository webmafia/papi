package openapi

import "iter"

type encoderContext struct {
	tags map[*Tag]struct{}
	refs map[string]*Ref
}

func newEncoderContext() *encoderContext {
	return &encoderContext{
		tags: make(map[*Tag]struct{}),
		refs: make(map[string]*Ref),
	}
}

func (ctx *encoderContext) addTag(tag *Tag) {
	ctx.tags[tag] = struct{}{}
}

func (ctx *encoderContext) addRef(ref *Ref) {
	ctx.refs[ref.Name] = ref
}

func (ctx *encoderContext) allRefs() iter.Seq2[int, *Ref] {
	return func(yield func(int, *Ref) bool) {
		var i int
		done := make(map[*Ref]struct{})

		for {
			var changed bool

			for _, ref := range ctx.refs {
				if _, ok := done[ref]; ok {
					continue
				}

				if !yield(i, ref) {
					return
				}

				done[ref] = struct{}{}
				changed = true
				i++
			}

			if !changed {
				break
			}
		}
	}
}
