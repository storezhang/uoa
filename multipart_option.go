package uoa

type multipartOption interface {
	applyMultipart(options *multipartOptions)
}

// NewMultipartOption 创建选项，因为option接口不对外暴露，如果用户想在外面创建option并赋值将无法完成，特意提供创建option的快捷方式
func NewMultipartOption(opts ...multipartOption) []multipartOption {
	return opts
}
