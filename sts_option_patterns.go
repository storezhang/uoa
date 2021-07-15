package uoa

var _ stsOption = (*stsOptionPatterns)(nil)

type stsOptionPatterns struct {
	patterns []string
}

// Patterns 配置多文件名
func Patterns(patterns ...string) *stsOptionPatterns {
	return &stsOptionPatterns{
		patterns: patterns,
	}
}

func (p *stsOptionPatterns) applySts(options *stsOptions) {
	options.patterns = p.patterns
}
