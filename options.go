package uoa

import (
	`time`

	`github.com/storezhang/gox`
)

var defaultOptions = &options{
	expired:   24 * time.Hour,
	separator: "/",
}

type options struct {
	// 通信端点
	endpoint string
	// 授权密钥
	secret gox.Secret
	// 过期时间
	expired time.Duration
	// 环境
	environment string
	// 分隔符
	separator string

	// 类型
	uoaType Type
}
