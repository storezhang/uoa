package uoa

var _ deleteOption = (*optionVersion)(nil)

type optionVersion struct {
	version string
}

// Version 配置版本
func Version(version string) *optionVersion {
	return &optionVersion{
		version: version,
	}
}

func (v *optionVersion) applyDelete(options *deleteOptions) {
	options.version = v.version
}
