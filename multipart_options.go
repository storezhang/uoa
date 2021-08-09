package uoa

type multipartOptions struct {
	*options

	objects []Object
}

func defaultMultipartOptions() *multipartOptions {
	return &multipartOptions{
		options: defaultOptions,
	}
}
