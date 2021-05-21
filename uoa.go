package uoa

import (
	`context`

	`github.com/storezhang/validatorx`
)

// Uoa 对象存储接口
type Uoa interface {
	// UploadUrl 上传地址
	UploadUrl(ctx context.Context, key Key, opts ...option) (uploadUrl string, err error)
	// DownloadUrl 下载地址
	DownloadUrl(ctx context.Context, key Key, filename string, opts ...option) (downloadUrl string, err error)
}

// New 创建适配器
func New(config Config) (uoa Uoa, err error) {
	var implementer Uoa

	if err = validatorx.Validate(config); nil != err {
		return
	}

	switch config.Type {
	case TypeCos:
		implementer, err = NewCos(config.Cos)
	}
	uoa = &uoaMaker{uoa: implementer}

	return
}
