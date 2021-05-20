package uoa

import (
	`context`

	`github.com/mcuadros/go-defaults`
)

// Uoa 对象存储接口
type Uoa interface {
	// UploadUrl 上传地址
	UploadUrl(ctx context.Context, key Key, opts ...option) (uploadUrl string, err error)
	// DownloadUrl 下载地址
	DownloadUrl(ctx context.Context, key Key, filename string, opts ...option) (url string, err error)
}

// New 创建适配器
func New(config Config) (uoa Uoa, err error) {
	// 处理默认值
	defaults.SetDefaults(&config)

	switch config.Type {
	case TypeCos:
		uoa, err = NewCos(config.Cos)
	}

	return
}
