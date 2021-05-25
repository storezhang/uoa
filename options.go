package uoa

import (
	`time`

	`github.com/storezhang/gox`
)

type options struct {
	// 通信端点
	endpoint string
	// 授权密钥
	secret gox.Secret
	// 过期时间
	expired time.Duration
	// 下载文件
	download bool
	// 打开文件
	inline bool
	// 文件类型
	contentType string
	// 环境
	environment string
	// 分隔符
	separator string

	// 类型
	uoaType Type
}

func defaultOptions() *options {
	return &options{
		expired:   24 * time.Hour,
		download:  false,
		inline:    true,
		separator: "/",
	}
}
