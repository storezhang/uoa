package uoa

type option interface {
	apply(options *options)
}
