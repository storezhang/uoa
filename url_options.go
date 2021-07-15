package uoa

import (
	`time`
)

type urlOptions struct {
	baseOptions

	// 下载文件
	download bool
	// 打开文件
	inline bool
	// 文件类型
	contentType string
}

func defaultUrlOptions() *urlOptions {
	return &urlOptions{
		baseOptions: baseOptions{
			expired:   24 * time.Hour,
			separator: "/",
		},
		download: false,
		inline:   true,
	}
}
