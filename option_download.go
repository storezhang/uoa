package uoa

var _ option = (*optionDownload)(nil)

type optionDownload struct {
	isDownload bool
}

// Download 配置应用名称
func Download() *optionDownload {
	return &optionDownload{
		isDownload: true,
	}
}

func (b *optionDownload) apply(options *options) {
	options.isDownload = b.isDownload
}
