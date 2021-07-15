package uoa

import (
	`time`
)

type stsOptions struct {
	baseOptions

	// 版本
	version string
	// 区域
	region string
}

func defaultStsOptions() *stsOptions {
	return &stsOptions{
		baseOptions: baseOptions{
			expired:   30 * time.Minute,
			separator: "/",
		},
		version: "2.0",
		region:  "ap-guangzhou",
	}
}
