package uoa

var _ urlOption = (*optionSeparator)(nil)

type optionSeparator struct {
	separator string
}

// Separator 配置分隔符
func Separator(separator string) *optionSeparator {
	return &optionSeparator{
		separator: separator,
	}
}

func (s *optionSeparator) apply(options *options) {
	options.separator = s.separator
}

func (s *optionSeparator) applyUrl(options *urlOptions) {
	options.separator = s.separator
}

func (s *optionSeparator) applyCredential(options *credentialsOptions) {
	options.separator = s.separator
}
