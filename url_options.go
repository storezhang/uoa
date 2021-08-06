package uoa

type urlOptions struct {
	*options

	// 下载文件
	download bool
	// 打开文件
	inline bool
	// 文件名
	filename string
	// 文件类型
	contentType string
	// 流类型
	streamType streamType
}

func defaultUrlOptions() *urlOptions {
	return &urlOptions{
		options: defaultOptions,

		download:   false,
		inline:     true,
		streamType: streamTypeDownstream,
	}
}
