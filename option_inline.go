package uoa

var _ option = (*optionInline)(nil)

type optionInline struct {
	inline bool
}

// Inline 配置应用名称
func Inline() *optionInline {
	return &optionInline{
		inline: true,
	}
}

func (b *optionInline) apply(options *options) {
	options.inline = b.inline
}
