package uoa

import (
	`context`
	`fmt`
	`net/url`
	`sync`
	`time`

	`github.com/aws/aws-sdk-go/aws`
	c `github.com/aws/aws-sdk-go/aws/credentials`
	`github.com/aws/aws-sdk-go/aws/session`
	`github.com/aws/aws-sdk-go/service/s3`
)

type _s3 struct {
	clientCache    sync.Map
	paramCiProcess string
	paramPm3u8     string
	paramExpires   string

	expiresMin float64
	expiresMax float64
}

func newS3() *_s3 {
	return &_s3{
		clientCache: sync.Map{},

		paramCiProcess: `ci-process`,
		paramPm3u8:     `pm3u8`,
		paramExpires:   `expires`,

		expiresMin: 3600,
		expiresMax: 43200,
	}
}

func newS3Client(options *options) (s *s3.S3, err error) {
	var sess *session.Session
	sess, err = session.NewSession(&aws.Config{
		Credentials: c.NewStaticCredentials(options.secret.Id, options.secret.Key, ""),
	})
	s = s3.New(sess)

	return
}

func (a *_s3) exist(ctx context.Context, bucket string, key string, options *options) (exist bool, err error) {
	var client *s3.S3
	client, err = newS3Client(options)
	if nil != err {
		return
	}

	if headRsp, headErr := client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); nil != headErr {
		exist = false
		err = headErr
	} else {
		exist = nil != headRsp
	}

	return
}

func (a *_s3) credentials(ctx context.Context, options *credentialsOptions, keys ...string) (credentials *credentialsBase, err error) {
	credential := c.NewStaticCredentials(options.secret.Id, options.secret.Key, "")

	var val c.Value
	val, err = credential.Get()
	if nil != err {
		return
	}

	expTime, _ := credential.ExpiresAt()
	credentials = &credentialsBase{
		Id:      val.AccessKeyID,
		Key:     val.SecretAccessKey,
		Token:   val.SessionToken,
		Expired: expTime,
	}

	return
}

func (a *_s3) url(ctx context.Context, bucket string, key string, options *urlOptions) (url *url.URL, err error) {
	var client *s3.S3
	client, err = newS3Client(options.options)
	if nil != err {
		return
	}

	switch options.streamType {
	case streamTypeUpstream:
		url, err = a.uploadUrl(ctx, client, bucket, key, options)
	case streamTypeDownstream:
		url, err = a.downloadUrl(ctx, client, bucket, key, options)
	default:
		url, err = a.downloadUrl(ctx, client, bucket, key, options)
	}

	return
}

func (a *_s3) initiateMultipart(ctx context.Context, key string, options *multipartOptions) (uploadId string, err error) {
	var client *s3.S3
	client, err = newS3Client(options.options)
	if nil != err {
		return
	}

	var res *s3.CreateMultipartUploadOutput
	if res, err = client.CreateMultipartUploadWithContext(ctx, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(options.bucket),
		Key:    aws.String(key),
	}); err != nil {
		return
	}
	uploadId = *res.UploadId

	return
}

func (a *_s3) completeMultipart(ctx context.Context, key string, uploadId string, objects []Object, options *multipartOptions) (err error) {
	var client *s3.S3
	var partNum int64
	client, err = newS3Client(options.options)
	if nil != err {
		return
	}

	parts := make([]*s3.CompletedPart, 0, len(objects))
	for _, object := range objects {
		partNum = int64(object.part)
		parts = append(parts, &s3.CompletedPart{
			ETag:       &object.etag,
			PartNumber: &partNum,
		})
	}
	partsUpload := &s3.CompletedMultipartUpload{
		Parts: parts,
	}
	_, err = client.CompleteMultipartUploadWithContext(ctx, &s3.CompleteMultipartUploadInput{
		Key:             aws.String(key),
		UploadId:        aws.String(uploadId),
		Bucket:          aws.String(options.bucket),
		MultipartUpload: partsUpload,
	})

	return
}

func (a *_s3) abortMultipart(ctx context.Context, key string, uploadId string, options *multipartOptions) (err error) {
	var client *s3.S3
	client, err = newS3Client(options.options)
	if nil != err {
		return
	}

	_, err = client.AbortMultipartUploadWithContext(ctx, &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(options.bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadId),
	})

	return
}

func (a *_s3) delete(ctx context.Context, key string, options *deleteOptions) (err error) {
	var client *s3.S3
	client, err = newS3Client(options.options)
	if nil != err {
		return
	}
	_, err = client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(options.bucket),
		Key:    aws.String(key),
	})

	return
}

func (a *_s3) uploadUrl(ctx context.Context, client *s3.S3, bucket string, key string, options *urlOptions) (url *url.URL, err error) {
	putOption := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	req, _ := client.PutObjectRequest(putOption)
	urlStr, _, err := req.PresignRequest(time.Duration(options.expired.Seconds()))
	if nil != err {
		return
	}
	if "" == urlStr {
		return
	}
	url, err = url.Parse(urlStr)

	return
}

func (a *_s3) downloadUrl(ctx context.Context, client *s3.S3, bucket string, key string, options *urlOptions) (url *url.URL, err error) {
	// 检查文件是否存在，文件不存在没必要往下继续执行
	var exist bool
	exist, err = a.exist(ctx, bucket, key, options.options)
	if !exist || nil != err {
		return
	}

	getOption := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	req, _ := client.GetObjectRequest(getOption)
	urlStr, _, err := req.PresignRequest(time.Duration(options.expired.Seconds()))
	if nil != err {
		return
	}
	if "" == urlStr {
		return
	}
	url, err = url.Parse(urlStr)

	// 解析私有M3u8存储文件
	if !options.pm3u8 {
		return
	}

	query := url.Query()
	query.Add(a.paramCiProcess, a.paramPm3u8)
	query.Add(a.paramExpires, a.expires(options.options))

	return
}

func (a *_s3) expires(options *options) string {
	expires := options.expired.Seconds()
	if expires < a.expiresMin {
		expires = a.expiresMin
	}
	if expires > a.expiresMax {
		expires = a.expiresMax
	}

	return fmt.Sprintf(`%f`, expires)
}
