package uoa

import (
	`context`
	`net/url`
	`strings`
)

// 内部接口封装
// 使用模板方法设计模式
type uoaTemplate struct {
	cos uoaInternal
}

func (t *uoaTemplate) UploadUrl(ctx context.Context, key Path, opts ...urlOption) (uploadUrl string, err error) {
	options := defaultDownloadOptions()
	for _, opt := range opts {
		opt.applyUrl(options)
	}

	fileKey := t.fileKey(key, options.environment, options.separator)
	var originalURL *url.URL
	switch options.uoaType {
	case TypeCos:
		originalURL, err = t.cos.uploadUrl(ctx, fileKey, options)
	}
	if nil != err {
		return
	}
	// 解决Golang JSON序列化时的HTML Escape
	uploadUrl = t.escape(originalURL)

	return
}

func (t *uoaTemplate) DownloadUrl(ctx context.Context, path Path, filename string, opts ...urlOption) (downloadUrl string, err error) {
	options := defaultDownloadOptions()
	for _, opt := range opts {
		opt.applyUrl(options)
	}

	fileKey := t.fileKey(path, options.environment, options.separator)
	var originalURL *url.URL
	switch options.uoaType {
	case TypeCos:
		originalURL, err = t.cos.downloadUrl(ctx, fileKey, filename, options)
	}
	if nil != err {
		return
	}
	// 解决Golang JSON序列化时的HTML Escape
	downloadUrl = t.escape(originalURL)

	return
}

func (t *uoaTemplate) fileKey(key Path, environment string, separator string) (fileKey string) {
	paths := key.Paths()
	if "" != environment {
		paths = append([]string{environment}, paths...)
	}
	fileKey = strings.Join(key.Paths(), separator)

	return
}

func (t *uoaTemplate) escape(originalURL *url.URL) (url string) {
	url = originalURL.String()
	url = strings.Replace(url, "\\u003c", "<", -1)
	url = strings.Replace(url, "\\u003e", ">", -1)
	url = strings.Replace(url, "\\u0026", "&", -1)

	return url
}
