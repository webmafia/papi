package openapi

import (
	"fmt"
	"iter"
)

type encoderContext struct {
	tags map[string]Tag
	refs map[string]*Ref
	auth bool
}

func newEncoderContext() *encoderContext {
	return &encoderContext{
		tags: make(map[string]Tag),
		refs: make(map[string]*Ref),
	}
}

func (ctx *encoderContext) addTag(tag Tag) {
	ctx.tags[tag.Name] = tag
}

func (ctx *encoderContext) addRef(ref *Ref) (err error) {
	if cur, ok := ctx.refs[ref.Name]; ok {
		if cur.Hash() != ref.Hash() {
			err = fmt.Errorf("reference name '%s' already exists", ref.Name)
		}

		return
	}

	ctx.refs[ref.Name] = ref
	return
}

func (ctx *encoderContext) allRefs() iter.Seq2[int, *Ref] {
	return func(yield func(int, *Ref) bool) {
		var i int
		done := make(map[string]struct{})

		for {
			var changed bool

			for _, ref := range ctx.refs {
				if _, ok := done[ref.Name]; ok {
					continue
				}

				if !yield(i, ref) {
					return
				}

				done[ref.Name] = struct{}{}
				changed = true
				i++
			}

			if !changed {
				break
			}
		}
	}
}
