package openapi

type encoderContext struct {
	tags map[*Tag]struct{}
	refs map[*Ref]struct{}
}

func newEncoderContext() *encoderContext {
	return &encoderContext{
		tags: make(map[*Tag]struct{}),
		refs: make(map[*Ref]struct{}),
	}
}

func (ctx *encoderContext) addTag(tag *Tag) {
	ctx.tags[tag] = struct{}{}
}

func (ctx *encoderContext) addRef(ref *Ref) {
	ctx.refs[ref] = struct{}{}
}
