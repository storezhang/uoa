package uoa

var _ urlOption = (*optionStream)(nil)

type optionStream struct {
	streamType streamType
}

// Upstream 配置上传
func Upstream() *optionStream {
	return &optionStream{streamType: streamTypeUpstream}
}

// Downstream 配置下载
func Downstream() *optionStream {
	return &optionStream{streamType: streamTypeDownstream}
}

func (e *optionStream) applyUrl(options *urlOptions) {
	options.streamType = e.streamType
}

func (e *optionStream) applyCredential(options *credentialsOptions) {
	options.streamType = e.streamType
}
