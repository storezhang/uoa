package uoa

var _ urlOption = (*optionContentType)(nil)

type optionContentType struct {
	contentType string
}

// ContentType 配置应用名称
func ContentType(contentType string) *optionContentType {
	return &optionContentType{
		contentType: contentType,
	}
}

func (ct *optionContentType) applyUrl(options *urlOptions) {
	options.contentType = ct.contentType
}
