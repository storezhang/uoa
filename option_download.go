package uoa

var _ urlOption = (*optionDownload)(nil)

type optionDownload struct {
	download bool
}

// Download 配置应用名称
func Download() *optionDownload {
	return &optionDownload{
		download: true,
	}
}

func (d *optionDownload) applyUrl(options *urlOptions) {
	options.download = d.download
}
