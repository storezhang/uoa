package uoa

var _ deleteOption = (*deleteOptionVersion)(nil)

type deleteOptionVersion struct {
	version string
}

// Version 配置版本
func Version(version string) *deleteOptionVersion {
	return &deleteOptionVersion{
		version: version,
	}
}

func (v *deleteOptionVersion) applyDelete(options *deleteOptions) {
	options.version = v.version
}
