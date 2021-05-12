package uoa

var _ option = (*optionContentType)(nil)

type optionContentType struct {
	contentType string
}

// ContentType 配置应用名称
func ContentType(contentType string) *optionContentType {
	return &optionContentType{
		contentType: contentType,
	}
}

func (b *optionContentType) apply(options *options) {
	options.contentType = b.contentType
}
