package uoa

// Credentials 授权
type Credentials struct {
	*credentialsBase

	// 连接地址
	Url string `json:"url" yaml:"url" xml:"url"`
	// 分隔符
	Separator string `json:"separator" yaml:"separator" xml:"separator"`
}
