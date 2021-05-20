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
		// 样式分隔符
		Separator string `default:"/" json:"separator" yaml:"separator" validate:"len=1"`
	}

	// Cos 腾讯云存储
	Cos struct {
		config CosConfig

		client *cos.Client
	}
)

// NewCos 创建腾讯云对象存储实现类
func NewCos(config CosConfig) (client *Cos, err error) {
	var bucketUrl *url.URL
	if bucketUrl, err = url.Parse(config.Url); nil != err {
		return
	}

	client = &Cos{
		config: config,

		client: cos.NewClient(&cos.BaseURL{BucketURL: bucketUrl}, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  config.Secret.Id,
				SecretKey: config.Secret.Key,
				// nolint:gosec
				Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
			},
		}),
	}

	return
}

func (c *Cos) UploadUrl(ctx context.Context, key Key, opts ...option) (uploadUrl string, err error) {
	appliedOptions := defaultOptions()
	for _, opt := range opts {
		opt.apply(&appliedOptions)
	}

	// 处理样式分隔符
	fileKey := strings.Join(key.Paths(), c.config.Separator)
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
		fileKey,
		c.config.Secret.Id, c.config.Secret.Key,
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

func (c *Cos) DownloadUrl(ctx context.Context, key Key, filename string, opts ...option) (downloadUrl string, err error) {
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

	// 处理样式分隔符
	fileKey := strings.Join(key.Paths(), c.config.Separator)
	// 检查文件是否存在，文件不存在没必要往下继续执行
	if headRsp, err = c.client.Object.Head(ctx, fileKey, nil); nil != err {
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
		fileKey,
		c.config.Secret.Id, c.config.Secret.Key,
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
