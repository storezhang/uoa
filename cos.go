package uoa

import (
	`context`
	`crypto/tls`
	`net/http`
	`net/url`
	`strings`

	`github.com/storezhang/gox`
	`github.com/tencentyun/cos-go-sdk-v5`
)

type (
	CosConfig struct {
		// 授权
		Secret gox.Secret `json:"secret" yaml:"secret" validate:"required"`
		// 存储桶地址
		Url string `json:"url" yaml:"url" validate:"required,url"`
	}

	// Cos 腾讯云存储
	Cos struct {
		secret gox.Secret

		client *cos.Client
	}
)

// NewCos 创建腾讯云对象存储实现类
func NewCos(secret gox.Secret, url *cos.BaseURL) *Cos {
	return &Cos{
		secret: secret,

		client: cos.NewClient(url, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  secret.Id,
				SecretKey: secret.Key,
				// nolint:gosec
				Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
			},
		}),
	}
}

func (c *Cos) UploadUrl(ctx context.Context, key string, opts ...option) (uploadUrl string, err error) {
	appliedOptions := defaultOptions()
	for _, opt := range opts {
		opt.apply(&appliedOptions)
	}

	var preassignedURL *url.URL
	putOptions := cos.ObjectPutHeaderOptions{
		XOptionHeader: &http.Header{
			"Access-Control-Expose-Headers": []string{"ETag"},
		},
	}
	// 获取预签名URL
	if preassignedURL, err = c.client.Object.GetPresignedURL(
		ctx,
		http.MethodPut,
		key,
		c.secret.Id, c.secret.Key,
		appliedOptions.expired,
		putOptions,
	); nil != err {
		return
	}
	uploadUrl = preassignedURL.String()
	// 解决Golang JSON序列化时的HTML Escape
	uploadUrl = c.escape(preassignedURL.String())

	return
}

func (c *Cos) DownloadUrl(ctx context.Context, key string, filename string, opts ...option) (downloadUrl string, err error) {
	appliedOptions := defaultOptions()
	for _, opt := range opts {
		opt.apply(&appliedOptions)
	}

	var (
		preassignedURL *url.URL
		getOptions     *cos.ObjectGetOptions
		headRsp        *cos.Response
		contentType    string
	)

	// 检查文件是否存在，文件不存在没必要往下继续执行
	if headRsp, err = c.client.Object.Head(ctx, key, nil); nil != err {
		if rspErr, ok := err.(*cos.ErrorResponse); ok && http.StatusNotFound == rspErr.Response.StatusCode {
			err = nil
		}

		return
	}

	if appliedOptions.isDownload {
		getOptions = &cos.ObjectGetOptions{
			ResponseContentDisposition: gox.ContentDisposition(filename, gox.ContentDispositionTypeAttachment),
		}
	} else if appliedOptions.isInline {
		if "" != appliedOptions.contentType {
			contentType = appliedOptions.contentType
		} else {
			contentType = headRsp.Header.Get(gox.HeaderContentType)
		}
		getOptions = &cos.ObjectGetOptions{
			ResponseContentDisposition: gox.ContentDisposition(filename, gox.ContentDispositionTypeInline),
			ResponseContentType:        contentType,
		}
	}

	// 获取预签名URL
	if preassignedURL, err = c.client.Object.GetPresignedURL(
		ctx,
		http.MethodGet,
		key,
		c.secret.Id, c.secret.Key,
		appliedOptions.expired,
		getOptions,
	); nil != err {
		return
	}
	// 解决Golang JSON序列化时的HTML Escape
	downloadUrl = c.escape(preassignedURL.String())

	return
}

func (c *Cos) escape(url string) string {
	url = strings.Replace(url, "\\u003c", "<", -1)
	url = strings.Replace(url, "\\u003e", ">", -1)
	url = strings.Replace(url, "\\u0026", "&", -1)

	return url
}
