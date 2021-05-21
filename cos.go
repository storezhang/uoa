package uoa

import (
	`context`
	`crypto/tls`
	`net/http`
	`net/url`
	`strings`

	`github.com/mcuadros/go-defaults`
	`github.com/storezhang/gox`
	`github.com/storezhang/validatorx`
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
		config CosConfig

		client *cos.Client
	}
)

// NewCos 创建腾讯云对象存储实现类
func NewCos(config CosConfig) (client *Cos, err error) {
	// 处理默认值
	defaults.SetDefaults(&config)
	if err = validatorx.Validate(config); nil != err {
		return
	}

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
		opt.apply(appliedOptions)
	}

	// 处理样式分隔符
	fileKey := strings.Join(key.Paths(), appliedOptions.separator)
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

	return
}

func (c *Cos) DownloadUrl(ctx context.Context, key Key, filename string, opts ...option) (downloadUrl string, err error) {
	appliedOptions := defaultOptions()
	for _, opt := range opts {
		opt.apply(appliedOptions)
	}

	var (
		preassignedURL *url.URL
		getOptions     *cos.ObjectGetOptions
		headRsp        *cos.Response
		contentType    string
	)

	// 处理样式分隔符
	fileKey := strings.Join(key.Paths(), appliedOptions.separator)
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
	downloadUrl = preassignedURL.String()

	return
}
