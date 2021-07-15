package uoa

type urlOption interface {
	applyUrl(options *urlOptions)
}

// NewUrlOptions 创建选项，因为option接口不对外暴露，如果用户想在外面创建option并赋值将无法完成，特意提供创建option的快捷方式
func NewUrlOptions(opts ...urlOption) []urlOption {
	return opts
}
