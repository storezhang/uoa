package uoa

var _ credentialsOption = (*optionPatterns)(nil)

type optionPatterns struct {
	patterns []string
}

// Patterns 配置多文件名
func Patterns(patterns ...string) *optionPatterns {
	return &optionPatterns{
		patterns: patterns,
	}
}

func (p *optionPatterns) applyCredential(options *credentialsOptions) {
	options.patterns = p.patterns
}
