package uoa

var _ urlOption = (*optionPm3u8)(nil)

type optionPm3u8 struct {
	pm3u8 bool
}

// Pm3u8 配置解析私有M3u8存储文件
func Pm3u8() *optionPm3u8 {
	return &optionPm3u8{
		pm3u8: true,
	}
}

func (p *optionPm3u8) applyUrl(options *urlOptions) {
	options.pm3u8 = p.pm3u8
}
