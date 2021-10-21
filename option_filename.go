package uoa

var _ urlOption = (*optionFilename)(nil)

type optionFilename struct {
	filename string
}

// Filename 配置文件名
func Filename(filename string) *optionFilename {
	return &optionFilename{
		filename: filename,
	}
}

func (f *optionFilename) applyUrl(options *urlOptions) {
	options.filename = f.filename
}
