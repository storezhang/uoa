package uoa

// Config 配置
type Config struct {
	// 类型
	Type Type `json:"type" yaml:"type" validate:"required,oneof=cos"`
	// 腾讯云对象存储
	Cos cosConfig `json:"tencentyun" yaml:"tencentyun" validate:"structonly"`
}
