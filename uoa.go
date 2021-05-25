package uoa

import (
	`context`
)

// Uoa 对象存储接口
type Uoa interface {
	// UploadUrl 上传地址
	UploadUrl(ctx context.Context, key Key, opts ...option) (uploadUrl string, err error)
	// DownloadUrl 下载地址
	DownloadUrl(ctx context.Context, key Key, filename string, opts ...option) (downloadUrl string, err error)
}

// New 创建适配器
func New() Uoa {
	return &uoaTemplate{
		cos: NewCos(),
	}
}
