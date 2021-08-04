package uoa

type multipartOptions struct {
	*options
}

func defaultMultipartOptions() *multipartOptions {
	return &multipartOptions{
		options: defaultOptions,
	}
}
