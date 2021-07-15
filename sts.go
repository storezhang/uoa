package uoa

// Sts 临时授权
type Sts struct {
	// 临时授权，相当于用户名
	Id string `json:"id" yaml:"id" xml:"id"`
	// 临时授权，相当于密码
	Key string `json:"key" yaml:"key" xml:"key"`
	// 临时授权
	Token string `json:"token" yaml:"token" xml:"token"`
}
