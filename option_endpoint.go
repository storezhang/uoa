package uoa

var _ option = (*optionEndpoint)(nil)

type optionEndpoint struct {
	// 通信端点
	endpoint string
}

// Endpoint 配置通信端点
func Endpoint(endpoint string) *optionEndpoint {
	return &optionEndpoint{
		endpoint: endpoint,
	}
}

// CosUrl 配置Cos地址
func CosUrl(url string) *optionEndpoint {
	return &optionEndpoint{
		endpoint: url,
	}
}

func (b *optionEndpoint) apply(options *options) {
	options.endpoint = b.endpoint
}
