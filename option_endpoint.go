package uoa

var _ urlOption = (*optionEndpoint)(nil)

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

func (e *optionEndpoint) applyUrl(options *urlOptions) {
	options.endpoint = e.endpoint
	options.uoaType = e.uoaType
}

func (e *optionEndpoint) applySts(options *stsOptions) {
	options.endpoint = e.endpoint
	options.uoaType = e.uoaType
}
