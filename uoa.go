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
func New(uoaType Type, validate *validatorx.Validate) (uoa Uoa, err error) {
	if err = validate.Var(uoaType, "required,oneof=cos"); nil != err {
		return
	}

	switch uoaType {
	case TypeCos:
		uoa = NewCos()
	}

	return
}
