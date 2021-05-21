package uoa

var _ option = (*optionSeparator)(nil)

type optionSeparator struct {
	separator string
}

// Separator 配置分隔符
func Separator(separator string) *optionSeparator {
	return &optionSeparator{
		separator: separator,
	}
}

func (b *optionSeparator) apply(options *options) {
	options.separator = b.separator
}
