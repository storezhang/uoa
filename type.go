package uoa

const (
	// TypeCos 腾讯云对象存储
	TypeCos Type = "cos"
	// S3 AWS云对象存储
	TypeS3 Type = "s3"
	// OBS 华为云对象存储
	TypeObs Type = "obs"
)

// Type 对象存储类型
type Type string
