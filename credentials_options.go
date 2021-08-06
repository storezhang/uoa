package uoa

type credentialsOptions struct {
	*options

	// 通信地址
	url string
	// 版本
	version string
	// 文件匹配
	patterns []string
	// 流类型
	streamType streamType
}

func defaultCredentialOptions() *credentialsOptions {
	return &credentialsOptions{
		options: defaultOptions,

		url:        "sts.tencentcloudapi.com",
		version:    "2.0",
		streamType: streamTypeUpstream,
	}
}
