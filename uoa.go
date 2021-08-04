package uoa

import (
	`context`
	`net/url`
	`sync`
)

// Uoa 对象存储接口
type Uoa interface {
	// Credentials 临时密钥
	Credentials(ctx context.Context, path Path, opts ...credentialsOption) (credentials *Credentials, err error)
	// Url 地址
	Url(ctx context.Context, path Path, filename string, opts ...urlOption) (url *url.URL, err error)
	// InitiateMultipartUpload 初始化分块上传
	InitiateMultipartUpload(ctx context.Context, path Path, opts ...multipartOption) (uploadId string, err error)
	// CompleteMultipartUpload 完成分块上传
	CompleteMultipartUpload(ctx context.Context, path Path, uploadId string, parts interface{}, opts ...multipartOption) (err error)
	// AbortMultipartUpload 终止分块上传
	AbortMultipartUpload(ctx context.Context, path Path, uploadId string, opts ...multipartOption) (err error)
	// Delete 删除
	Delete(ctx context.Context, path Path, opts ...deleteOption) (err error)
}

// New 创建适配器
func New(opts ...option) Uoa {
	for _, opt := range opts {
		opt.apply(defaultOptions)
	}

	return &uoaTemplate{
		cos: &cosInternal{clientCache: sync.Map{}},
	}
}
