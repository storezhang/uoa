package uoa

var _ urlOption = (*urlOptionDownload)(nil)

type urlOptionDownload struct {
	download bool
}

// Download 配置应用名称
func Download() *urlOptionDownload {
	return &urlOptionDownload{
		download: true,
	}
}

func (d *urlOptionDownload) applyUrl(options *urlOptions) {
	options.download = d.download
}
