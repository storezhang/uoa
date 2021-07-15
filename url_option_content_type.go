package uoa

var _ urlOption = (*urlOptionContentType)(nil)

type urlOptionContentType struct {
	contentType string
}

// ContentType 配置应用名称
func ContentType(contentType string) *urlOptionContentType {
	return &urlOptionContentType{
		contentType: contentType,
	}
}

func (ct *urlOptionContentType) applyUrl(options *urlOptions) {
	options.contentType = ct.contentType
}
