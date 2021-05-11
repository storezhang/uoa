package uoa

import (
	`context`
	`net/url`

	`github.com/tencentyun/cos-go-sdk-v5`
)

// Uoa 对象存储接口
type Uoa interface {
	// UploadUrl 上传地址
	UploadUrl(ctx context.Context, key string, opts ...option) (uploadUrl string, err error)
	// DownloadUrl 下载地址
	DownloadUrl(ctx context.Context, key string, filename string, opts ...option) (url string, err error)
}

func New(config Config) (uoa Uoa, err error) {
	switch config.Type {
	case TypeCos:
		var bucketUrl *url.URL
		if bucketUrl, err = url.Parse(config.Cos.Url); nil != err {
			return
		}
		uoa = NewCos(config.Cos.Secret, &cos.BaseURL{BucketURL: bucketUrl})
	}

	return
}
