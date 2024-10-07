package openapi

type encoderContext struct {
	tags map[*Tag]struct{}
}

func newEncoderContext() *encoderContext {
	return &encoderContext{
		tags: make(map[*Tag]struct{}),
	}
}

func (ctx *encoderContext) addTag(tag *Tag) {
	ctx.tags[tag] = struct{}{}
}
