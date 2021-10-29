package uoa

import (
	`context`
	`net/url`
	`strings`
	`sync`
)

type obs struct {
	clientCache sync.Map
}

func newObs() *obs {
	return &obs{
		clientCache: sync.Map{},
	}
}

func newObsClient(accessKey, securityKey, endPoint string) (*obsClient, error) {
	return NewObsClient(accessKey, securityKey, endPoint)
}

func (o *obs) exist(ctx context.Context, key string, options *options) (exist bool, err error) {
	var (
		client *obsClient
		output *baseModel
	)

	client, err = newObsClient(options.secret.Id, options.secret.Key, options.endpoint)
	if nil != err {
		return
	}

	_, _, bucket := o.parseRegionAndBucketName(options.endpoint)
	output, err = client.headObject(&headObjectInput{
		bucket: bucket,
		key:    key,
	})
	if nil != err || output == nil {
		exist = false
		return
	}
	if output.statusCode != 200 {
		exist = false
	} else {
		exist = true
	}

	return
}

func (o *obs) credentials(ctx context.Context, options *credentialsOptions, keys ...string) (credentials *credentialsBase, err error) {
	var client *obsClient

	client, err = newObsClient(options.secret.Id, options.secret.Key, options.endpoint)
	if nil != err {
		return
	}

	securityHolder := client.getSecurity()
	credentials = &credentialsBase{
		Id:    securityHolder.accessKey,
		Key:   securityHolder.securityKey,
		Token: securityHolder.securityToken,
	}

	return
}

func (o *obs) url(ctx context.Context, key string, options *urlOptions) (url *url.URL, err error) {
	var (
		client *obsClient
		method HttpMethodType
	)

	client, err = newObsClient(options.secret.Id, options.secret.Key, options.endpoint)
	if nil != err {
		return
	}

	switch options.streamType {
	case streamTypeUpstream:
		method = HttpMethodPost
	case streamTypeDownstream:
		method = HttpMethodGet
	default:
		method = HttpMethodGet
	}

	_, _, bucket := o.parseRegionAndBucketName(options.endpoint)
	output, _err := client.createSignedUrl(&createSignedUrlInput{
		method:  method,
		bucket:  bucket,
		key:     key,
		expires: 0,
	})
	if nil != _err {
		err = _err
		return
	}
	urlStr := output.signedUrl
	if "" == urlStr {
		return
	}
	url, err = url.Parse(urlStr)

	return
}

func (o *obs) initiateMultipart(ctx context.Context, key string, options *multipartOptions) (uploadId string, err error) {
	var (
		client *obsClient
		output *initiateMultipartUploadOutput
	)

	client, err = newObsClient(options.secret.Id, options.secret.Key, options.endpoint)
	if nil != err {
		return
	}

	_, _, bucket := o.parseRegionAndBucketName(options.endpoint)
	objectOperationInput := &objectOperationInput{
		bucket: bucket,
		key:    key,
	}
	output, err = client.initiateMultipartUpload(&initiateMultipartUploadInput{
		objectOperationInput: *objectOperationInput,
	})

	if err != nil {
		return
	}
	uploadId = output.uploadId

	return
}

func (o *obs) completeMultipart(ctx context.Context, key string, uploadId string, objects []Object, options *multipartOptions) (err error) {
	var client *obsClient

	client, err = newObsClient(options.secret.Id, options.secret.Key, options.endpoint)
	if nil != err {
		return
	}

	parts := make([]part, 0, len(objects))
	for _, object := range objects {
		parts = append(parts, part{
			eTag:       object.etag,
			partNumber: int(object.size),
			size:       object.size,
		})
	}

	_, _, bucket := o.parseRegionAndBucketName(options.endpoint)
	input := &completeMultipartUploadInput{
		bucket:   bucket,
		key:      key,
		uploadId: uploadId,
		parts:    parts,
	}
	_, err = client.completeMultipartUpload(input)

	return
}

func (o *obs) abortMultipart(ctx context.Context, key string, uploadId string, options *multipartOptions) (err error) {
	var client *obsClient

	client, err = newObsClient(options.secret.Id, options.secret.Key, options.endpoint)
	if nil != err {
		return
	}

	_, _, bucket := o.parseRegionAndBucketName(options.endpoint)
	input := &abortMultipartUploadInput{
		bucket:   bucket,
		key:      key,
		uploadId: uploadId,
	}
	_, err = client.abortMultipartUpload(input)

	return
}

func (o *obs) delete(ctx context.Context, key string, options *deleteOptions) (err error) {
	var client *obsClient

	client, err = newObsClient(options.secret.Id, options.secret.Key, options.endpoint)
	if nil != err {
		return
	}

	_, _, bucket := o.parseRegionAndBucketName(options.endpoint)
	input := &deleteObjectInput{
		bucket: bucket,
		key:    key,
	}
	_, err = client.deleteObject(input)

	return
}

func (o *obs) parseRegionAndBucketName(endpoint string) (region string, appId string, bucketName string) {
	endpoint = strings.ReplaceAll(endpoint, `https://`, ``)
	urls := strings.Split(endpoint, `.`)
	region = urls[2]
	bucketName = urls[0]
	appId = strings.Split(urls[0], `-`)[1]

	return
}