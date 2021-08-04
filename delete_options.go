package uoa

type deleteOptions struct {
	*options

	// 版本
	version string
}

func defaultDeleteOptions() *deleteOptions {
	return &deleteOptions{
		options: defaultOptions,
	}
}
