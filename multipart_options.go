package uoa

type multipartOptions struct {
	*options

	objects []Object

	// 桶名称，非必填参数，针对S3需传入
	bucket string
}

func defaultMultipartOptions() *multipartOptions {
	return &multipartOptions{
		options: defaultOptions,
	}
}
