package uoa

var _ option = (*optionEndpoint)(nil)

type optionEndpoint struct {
	// 通信端点
	endpoint string
	// 类型
	uoaType Type
}

// Endpoint 配置通信端点
func Endpoint(endpoint string, uoaType Type) *optionEndpoint {
	return &optionEndpoint{
		endpoint: endpoint,
		uoaType:  uoaType,
	}
}

// CosUrl 配置Cos地址
func CosUrl(url string) *optionEndpoint {
	return Endpoint(url, TypeCos)
}

func (b *optionEndpoint) apply(options *options) {
	options.endpoint = b.endpoint
	options.uoaType = b.uoaType
}
