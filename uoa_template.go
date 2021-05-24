package uoa

import (
	`context`
	`net/url`
	`strings`
)

// 内部接口封装
// 使用模板方法设计模式
type uoaTemplate struct {
	implementer uoaInternal
}

func (t *uoaTemplate) UploadUrl(ctx context.Context, key Key, opts ...option) (uploadUrl string, err error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt.apply(options)
	}

	fileKey := t.fileKey(key, options.environment, options.separator)
	var originalURL *url.URL
	if originalURL, err = t.implementer.uploadUrl(ctx, fileKey, options); nil != err {
		return
	}
	// 解决Golang JSON序列化时的HTML Escape
	uploadUrl = t.escape(originalURL)

	return
}

func (t *uoaTemplate) DownloadUrl(ctx context.Context, key Key, filename string, opts ...option) (downloadUrl string, err error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt.apply(options)
	}

	var originalURL *url.URL
	fileKey := t.fileKey(key, options.environment, options.separator)
	if originalURL, err = t.implementer.downloadUrl(ctx, fileKey, filename, options); nil != err {
		return
	}
	// 解决Golang JSON序列化时的HTML Escape
	downloadUrl = t.escape(originalURL)

	return
}

func (t *uoaTemplate) fileKey(key Key, environment string, separator string) (fileKey string) {
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
