package uoa

var _ urlOption = (*optionInline)(nil)

type optionInline struct {
	inline bool
}

// Inline 配置应用名称
func Inline() *optionInline {
	return &optionInline{
		inline: true,
	}
}

func (i *optionInline) applyUrl(options *urlOptions) {
	options.inline = i.inline
}
