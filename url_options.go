package uoa

type urlOptions struct {
	*options

	// 下载文件
	download bool
	// 打开文件
	inline bool
	// 文件类型
	contentType string
}

func defaultUrlOptions() *urlOptions {
	return &urlOptions{
		options: defaultOptions,

		download: false,
		inline:   true,
	}
}
