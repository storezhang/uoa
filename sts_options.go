package uoa

import (
	`time`
)

type stsOptions struct {
	baseOptions

	// 通信地址
	url string
	// 版本
	version string

	// 文件匹配
	patterns []string
}

func defaultStsOptions() *stsOptions {
	return &stsOptions{
		baseOptions: baseOptions{
			expired:   30 * time.Minute,
			separator: "/",
		},
		url:     "sts.tencentcloudapi.com",
		version: "2.0",
	}
}
