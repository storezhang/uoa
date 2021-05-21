package uoa

import (
	`context`
	`strings`
)

// 内部接口封装
// 使用模板方法设计模式
type uoaTemplate struct {
	uoa Uoa
}

func (t *uoaTemplate) UploadUrl(ctx context.Context, key Key, opts ...option) (uploadUrl string, err error) {
	if uploadUrl, err = t.uoa.UploadUrl(ctx, key, opts...); nil != err {
		return
	}
	// 解决Golang JSON序列化时的HTML Escape
	uploadUrl = t.escape(uploadUrl)

	return
}

func (t *uoaTemplate) DownloadUrl(ctx context.Context, key Key, filename string, opts ...option) (downloadUrl string, err error) {
	appliedOptions := defaultOptions()
	for _, opt := range opts {
		opt.apply(appliedOptions)
	}

	keyTemplate := &keyTemplate{key: key, environment: appliedOptions.environment}
	if downloadUrl, err = t.uoa.DownloadUrl(ctx, keyTemplate, filename, opts...); nil != err {
		return
	}
	// 解决Golang JSON序列化时的HTML Escape
	downloadUrl = t.escape(downloadUrl)

	return
}

func (t *uoaTemplate) escape(url string) string {
	url = strings.Replace(url, "\\u003c", "<", -1)
	url = strings.Replace(url, "\\u003e", ">", -1)
	url = strings.Replace(url, "\\u0026", "&", -1)

	return url
}
