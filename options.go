package uoa

import (
	`time`
)

type options struct {
	// 过期时间
	expired time.Duration
	// 下载文件
	isDownload bool
	// 打开文件
	isInline bool
	// 文件类型
	contentType string
	// 环境
	environment string
	// 分隔符
	separator string
}

func defaultOptions() *options {
	return &options{
		expired:    24 * time.Hour,
		isDownload: false,
		isInline:   true,
		separator:  "/",
	}
}
