package uoa

var _ urlOption = (*urlOptionInline)(nil)

type urlOptionInline struct {
	inline bool
}

// Inline 配置应用名称
func Inline() *urlOptionInline {
	return &urlOptionInline{
		inline: true,
	}
}

func (i *urlOptionInline) applyUrl(options *urlOptions) {
	options.inline = i.inline
}
