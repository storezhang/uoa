package uoa

import (
	`context`
	`net/url`
	`strings`
	`sync`
	`time`

	`github.com/aws/aws-sdk-go/aws`
	c `github.com/aws/aws-sdk-go/aws/credentials`
	`github.com/aws/aws-sdk-go/aws/session`
	`github.com/aws/aws-sdk-go/service/s3`
)

type _s3 struct {
	clientCache    sync.Map
}

func newS3() *_s3 {
	return &_s3{
		clientCache: sync.Map{},
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

func (a *_s3) exist(ctx context.Context, key string, options *options) (exist bool, err error) {
	var client *s3.S3
	client, err = newS3Client(options)
	if nil != err {
		return
	}

	_, _, bucket := a.parseRegionAndBucket(options.endpoint)
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

func (a *_s3) url(ctx context.Context, key string, options *urlOptions) (url *url.URL, err error) {
	var client *s3.S3
	client, err = newS3Client(options.options)
	if nil != err {
		return
	}

	_, _, bucket := a.parseRegionAndBucket(options.endpoint)
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

	_, _, bucket := a.parseRegionAndBucket(options.endpoint)
	var res *s3.CreateMultipartUploadOutput
	if res, err = client.CreateMultipartUploadWithContext(ctx, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
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
	_, _, bucket := a.parseRegionAndBucket(options.endpoint)
	_, err = client.CompleteMultipartUploadWithContext(ctx, &s3.CompleteMultipartUploadInput{
		Key:             aws.String(key),
		UploadId:        aws.String(uploadId),
		Bucket:          aws.String(bucket),
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

	_, _, bucket := a.parseRegionAndBucket(options.endpoint)
	_, err = client.AbortMultipartUploadWithContext(ctx, &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(bucket),
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
	_, _, bucket := a.parseRegionAndBucket(options.endpoint)
	_, err = client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
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
	exist, err = a.exist(ctx, key, options.options)
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

	return
}

func (a *_s3) parseRegionAndBucket(endpoint string) (region string, appId string, bucketName string) {
	endpoint = strings.ReplaceAll(endpoint, `https://`, ``)
	urls := strings.Split(endpoint, `.`)
	region = urls[2]
	bucketName = urls[0]
	appId = strings.Split(urls[0], `-`)[1]

	return
}