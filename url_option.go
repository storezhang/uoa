package uoa

type (
	urlOption interface {
		applyUrl(options *urlOptions)
	}

	urlOptions struct {
		*options

		// 下载文件
		download bool
		// 打开文件
		inline bool
		// 是否处理M3u8私有存储
		pm3u8 bool
		// 文件名
		filename string
		// 文件类型
		contentType string
		// 流类型
		streamType streamType
	}
)

// NewUrlOptions 创建选项，因为option接口不对外暴露，如果用户想在外面创建option并赋值将无法完成，特意提供创建option的快捷方式
func NewUrlOptions(opts ...urlOption) []urlOption {
	return opts
}

func defaultUrlOptions() *urlOptions {
	return &urlOptions{
		options: defaultOptions,

		download:   false,
		inline:     true,
		pm3u8:      false,
		streamType: streamTypeDownstream,
	}
}
