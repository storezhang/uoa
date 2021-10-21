package uoa

type (
	credentialsOption interface {
		applyCredential(options *credentialsOptions)
	}

	credentialsOptions struct {
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
)

// NewCredentialsOptions 创建选项，因为option接口不对外暴露，如果用户想在外面创建option并赋值将无法完成，特意提供创建option的快捷方式
func NewCredentialsOptions(opts ...credentialsOption) []credentialsOption {
	return opts
}

func defaultCredentialOptions() *credentialsOptions {
	return &credentialsOptions{
		options: defaultOptions,

		url:        "sts.tencentcloudapi.com",
		version:    "2.0",
		streamType: streamTypeUpstream,
	}
}
