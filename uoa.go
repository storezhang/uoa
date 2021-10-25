package uoa

import (
	`context`
	`net/url`
)

// Uoa 对象存储接口
type Uoa interface {
	// Exist 检查文件是否存在
	Exist(ctx context.Context, bucket string, path Path, opts ...option) (exist bool, err error)

	// Credentials 临时密钥
	Credentials(ctx context.Context, path Path, opts ...credentialsOption) (credentials *Credentials, err error)

	// Url 地址
	Url(ctx context.Context, bucket string, path Path, opts ...urlOption) (url *url.URL, err error)

	// InitiateMultipart 初始化分块上传
	InitiateMultipart(ctx context.Context, path Path, opts ...multipartOption) (uploadId string, err error)

	// CompleteMultipart 完成分块上传
	CompleteMultipart(ctx context.Context, path Path, uploadId string, objects []Object, opts ...multipartOption) (err error)

	// AbortMultipart 终止分块上传
	AbortMultipart(ctx context.Context, path Path, uploadId string, opts ...multipartOption) (err error)

	// Delete 删除
	Delete(ctx context.Context, path Path, opts ...deleteOption) (err error)
}

// New 创建适配器
func New(opts ...option) Uoa {
	for _, opt := range opts {
		opt.apply(defaultOptions)
	}

	return &template{
		cos: newCos(),
		s3:  newS3(),
	}
}
