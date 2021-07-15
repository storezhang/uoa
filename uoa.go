package uoa

import (
	`context`
)

// Uoa 对象存储接口
type Uoa interface {
	// Sts 临时密钥
	Sts(ctx context.Context, path Path, opts ...stsOption) (sts Sts, err error)
	// Url 地址
	Url(ctx context.Context, path Path, filename string, opts ...urlOption) (url string, err error)
}

// New 创建适配器
func New() Uoa {
	return &uoaTemplate{
		cos: NewCos(),
	}
}
