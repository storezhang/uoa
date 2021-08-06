package uoa

type multipartOptions struct {
	*options

	objects []object
}

func defaultMultipartOptions() *multipartOptions {
	return &multipartOptions{
		options: defaultOptions,
	}
}
