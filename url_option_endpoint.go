package uoa

var _ urlOption = (*urlOptionEndpoint)(nil)

type urlOptionEndpoint struct {
	// 通信端点
	endpoint string
	// 类型
	uoaType Type
}

// Endpoint 配置通信端点
func Endpoint(endpoint string, uoaType Type) *urlOptionEndpoint {
	return &urlOptionEndpoint{
		endpoint: endpoint,
		uoaType:  uoaType,
	}
}

// CosUrl 配置Cos地址
func CosUrl(url string) *urlOptionEndpoint {
	return Endpoint(url, TypeCos)
}

func (e *urlOptionEndpoint) applyUrl(options *urlOptions) {
	options.endpoint = e.endpoint
	options.uoaType = e.uoaType
}
