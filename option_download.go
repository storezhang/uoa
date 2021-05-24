package uoa

var _ option = (*optionDownload)(nil)

type optionDownload struct {
	download bool
}

// Download 配置应用名称
func Download() *optionDownload {
	return &optionDownload{
		download: true,
	}
}

func (b *optionDownload) apply(options *options) {
	options.download = b.download
}
