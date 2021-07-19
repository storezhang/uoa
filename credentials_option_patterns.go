package uoa

var _ credentialsOption = (*credentialsOptionPatterns)(nil)

type credentialsOptionPatterns struct {
	patterns []string
}

// Patterns 配置多文件名
func Patterns(patterns ...string) *credentialsOptionPatterns {
	return &credentialsOptionPatterns{
		patterns: patterns,
	}
}

func (p *credentialsOptionPatterns) applyCredential(options *credentialsOptions) {
	options.patterns = p.patterns
}
