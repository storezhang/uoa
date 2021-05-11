package uoa

var _ option = (*optionInline)(nil)

type optionInline struct {
	isInline bool
}

// Inline 配置应用名称
func Inline() *optionInline {
	return &optionInline{
		isInline: true,
	}
}

func (b *optionInline) apply(options *options) {
	options.isInline = b.isInline
}
