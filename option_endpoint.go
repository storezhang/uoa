package uoa

var _ urlOption = (*optionEndpoint)(nil)

type optionEndpoint struct {
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
	return Endpoint(url)
}

func (e *optionEndpoint) apply(options *options) {
	options.endpoint = e.endpoint
}

func (e *optionEndpoint) applyUrl(options *urlOptions) {
	options.endpoint = e.endpoint
}

func (e *optionEndpoint) applyCredential(options *credentialsOptions) {
	options.endpoint = e.endpoint
}
