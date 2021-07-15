package uoa

import (
	`context`
	`crypto/tls`
	`fmt`
	`net/http`
	`net/url`
	`sync`

	`github.com/storezhang/gox`
	`github.com/tencentyun/cos-go-sdk-v5`
)

// Cos 腾讯云存储
type Cos struct {
	clientCache sync.Map

	template uoaTemplate
}

// NewCos 创建腾讯云对象存储实现类
func NewCos() (cos *Cos) {
	cos = &Cos{
		clientCache: sync.Map{},
	}
	cos.template = uoaTemplate{cos: cos}

	return
}

func (c *Cos) UploadUrl(ctx context.Context, key Path, opts ...urlOption) (uploadUrl string, err error) {
	return c.template.UploadUrl(ctx, key, opts...)
}

func (c *Cos) DownloadUrl(ctx context.Context, key Path, filename string, opts ...urlOption) (downloadUrl string, err error) {
	return c.template.DownloadUrl(ctx, key, filename, opts...)
}

func (c *Cos) sts(ctx context.Context, path Path, opts ...stsOption) (sts Sts, err error) {

	return
}

func (c *Cos) downloadUrl(ctx context.Context, key string, filename string, options *urlOptions) (downloadUrl *url.URL, err error) {
	var (
		client      *cos.Client
		getOptions  *cos.ObjectGetOptions
		headRsp     *cos.Response
		contentType string
	)

	if client, err = c.getClient(options.endpoint, options.secret); nil != err {
		return
	}
	// 检查文件是否存在，文件不存在没必要往下继续执行
	if headRsp, err = client.Object.Head(ctx, key, nil); nil != err {
		if rspErr, ok := err.(*cos.ErrorResponse); ok && http.StatusNotFound == rspErr.Response.StatusCode {
			err = nil
		}

		return
	}

	if options.download {
		getOptions = &cos.ObjectGetOptions{
			ResponseContentDisposition: gox.ContentDisposition(filename, gox.ContentDispositionTypeAttachment),
		}
	} else if options.inline {
		if "" != options.contentType {
			contentType = options.contentType
		} else {
			contentType = headRsp.Header.Get(gox.HeaderContentType)
		}
		getOptions = &cos.ObjectGetOptions{
			ResponseContentDisposition: gox.ContentDisposition(filename, gox.ContentDispositionTypeInline),
			ResponseContentType:        contentType,
		}
	}

	// 获取预签名URL
	if downloadUrl, err = client.Object.GetPresignedURL(
		ctx,
		http.MethodGet,
		key,
		options.secret.Id, options.secret.Key,
		options.expired,
		getOptions,
	); nil != err {
		return
	}

	return
}

func (c *Cos) getClient(baseUrl string, secret gox.Secret) (client *cos.Client, err error) {
	var (
		cache interface{}
		ok    bool
	)

	key := fmt.Sprintf("%s-%s", baseUrl, secret.Id)
	if cache, ok = c.clientCache.Load(key); ok {
		client = cache.(*cos.Client)

		return
	}

	var bucketUrl *url.URL
	if bucketUrl, err = url.Parse(baseUrl); nil != err {
		return
	}

	client = cos.NewClient(&cos.BaseURL{BucketURL: bucketUrl}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secret.Id,
			SecretKey: secret.Key,
			// nolint:gosec
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		},
	})
	c.clientCache.Store(key, client)

	return
}
