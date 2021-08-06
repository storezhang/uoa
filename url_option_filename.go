package uoa

var _ urlOption = (*urlOptionFilename)(nil)

type urlOptionFilename struct {
	filename string
}

// Filename 配置文件名
func Filename(filename string) *urlOptionFilename {
	return &urlOptionFilename{
		filename: filename,
	}
}

func (f *urlOptionFilename) applyUrl(options *urlOptions) {
	options.filename = f.filename
}
