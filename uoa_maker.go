package uoa

import (
	`context`
	`strings`
)

// 内部接口封装
// 使用模板方法设计模式
type uoaMaker struct {
	uoa Uoa
}

func (u *uoaMaker) UploadUrl(ctx context.Context, key Key, opts ...option) (uploadUrl string, err error) {
	if uploadUrl, err = u.uoa.UploadUrl(ctx, key, opts...); nil != err {
		return
	}
	// 解决Golang JSON序列化时的HTML Escape
	uploadUrl = u.escape(uploadUrl)

	return
}

func (u *uoaMaker) DownloadUrl(ctx context.Context, key Key, filename string, opts ...option) (downloadUrl string, err error) {
	appliedOptions := defaultOptions()
	for _, opt := range opts {
		opt.apply(&appliedOptions)
	}

	keyMaker := &keyMaker{key: key, environment: appliedOptions.environment}
	if downloadUrl, err = u.uoa.DownloadUrl(ctx, keyMaker, filename, opts...); nil != err {
		return
	}
	// 解决Golang JSON序列化时的HTML Escape
	downloadUrl = u.escape(downloadUrl)

	return
}

func (u *uoaMaker) escape(url string) string {
	url = strings.Replace(url, "\\u003c", "<", -1)
	url = strings.Replace(url, "\\u003e", ">", -1)
	url = strings.Replace(url, "\\u0026", "&", -1)

	return url
}
