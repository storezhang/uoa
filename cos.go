package uoa

import (
	`context`
	`crypto/tls`
	`encoding/json`
	`fmt`
	`net/http`
	`net/url`
	`strings`
	`sync`
	`time`

	`github.com/storezhang/gox`
	`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common`
	`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile`
	`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813`
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

func (c *Cos) Sts(ctx context.Context, path Path, opts ...stsOption) (sts Sts, err error) {
	return c.template.Sts(ctx, path, opts...)
}

func (c *Cos) DownloadUrl(ctx context.Context, path Path, filename string, opts ...urlOption) (downloadUrl string, err error) {
	return c.template.DownloadUrl(ctx, path, filename, opts...)
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

func (c *Cos) sts(_ context.Context, key string, options *stsOptions) (sts Sts, err error) {
	actions := []string{
		// 简单上传
		"name/cos:PutObject",
		// 表单上传、小程序上传
		"name/cos:PostObject",
		// 分片上传
		"name/cos:InitiateMultipartUpload",
		"name/cos:ListMultipartUploads",
		"name/cos:ListParts",
		"name/cos:UploadPart",
		"name/cos:CompleteMultipartUpload",
	}
	region, appId, bucketName := c.parse(options.endpoint)
	policy := cosPolicy{
		Version: options.version,
		Statements: []cosStatement{
			{
				Actions: actions,
				Effect:  "allow",
				Resource: []string{
					fmt.Sprintf("qcs::cos:%s:uid/%s:%s/%s", region, appId, bucketName, key),
				},
			},
		},
	}

	var policyBytes []byte
	if policyBytes, err = json.Marshal(policy); nil != err {
		return
	}

	credential := common.NewCredential(options.secret.Id, options.secret.Key)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = options.url
	client, _ := v20180813.NewClient(credential, region, cpf)

	req := v20180813.NewGetFederationTokenRequest()
	req.Name = common.StringPtr("cos-sts-go")
	req.Policy = common.StringPtr(string(policyBytes))
	req.DurationSeconds = common.Uint64Ptr(uint64(options.expired / time.Second))

	var rsp *v20180813.GetFederationTokenResponse
	if rsp, err = client.GetFederationToken(req); nil != err {
		return
	}

	sts = Sts{
		Id:      *rsp.Response.Credentials.TmpSecretId,
		Key:     *rsp.Response.Credentials.TmpSecretKey,
		Token:   *rsp.Response.Credentials.Token,
		Expired: time.Unix(int64(*rsp.Response.ExpiredTime), 0),
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

func (c *Cos) parse(endpoint string) (region string, appId string, bucketName string) {
	endpoint = strings.ReplaceAll(endpoint, "https://", "")
	urls := strings.Split(endpoint, ".")
	region = urls[2]
	bucketName = urls[0]
	appId = strings.Split(urls[0], "-")[1]

	return
}
