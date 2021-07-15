package uoa

import (
	`context`
)

// Uoa 对象存储接口
type Uoa interface {
	// Sts 临时密钥
	Sts(ctx context.Context, path Path, opts ...stsOption) (sts Sts, err error)
	// DownloadUrl 下载地址
	DownloadUrl(ctx context.Context, path Path, filename string, opts ...urlOption) (downloadUrl string, err error)
}

// New 创建适配器
func New() Uoa {
	return &uoaTemplate{
		cos: NewCos(),
	}
}
