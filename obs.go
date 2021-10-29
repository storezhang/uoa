package uoa

import (
	`context`
	`net/url`
	`sync`
)

type _obs struct {
	clientCache    sync.Map
	paramCiProcess string
	paramPm3u8     string
	paramExpires   string

	expiresMin float64
	expiresMax float64
}

func newObs() *_obs {
	return &_obs{
		clientCache: sync.Map{},

		paramCiProcess: `ci-process`,
		paramPm3u8:     `pm3u8`,
		paramExpires:   `expires`,

		expiresMin: 3600,
		expiresMax: 43200,
	}
}

func newObsClient(accessKey, securityKey, endPoint string) (*ObsClient, error) {
	return NewObsClient(accessKey, securityKey, endPoint)
}

func (_ *_obs) exist(ctx context.Context, bucket string, key string, options *options) (exist bool, err error) {
	var client *ObsClient
	var output *BaseModel
	client, err = newObsClient(options.secret.Id, options.secret.Key, options.endpoint)
	if nil != err {
		return
	}
	output, err = client.HeadObject(&HeadObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	if nil != err || output == nil {
		exist = false
		return
	}
	if output.StatusCode != 200 {
		exist = false
	} else {
		exist = true
	}

	return
}

func (_ *_obs) credentials(ctx context.Context, options *credentialsOptions, keys ...string) (credentials *credentialsBase, err error) {
	var client *ObsClient

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

func (_ *_obs) url(ctx context.Context, bucket string, key string, options *urlOptions) (url *url.URL, err error) {
	var client *ObsClient
	var method HttpMethodType

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

	output, _err := client.CreateSignedUrl(&CreateSignedUrlInput{
		Method:  method,
		Bucket:  bucket,
		Key:     key,
		Expires: 0,
	})
	if nil != _err {
		err = _err
		return
	}
	urlStr := output.SignedUrl
	if "" == urlStr {
		return
	}
	url, err = url.Parse(urlStr)

	return
}

func (_ *_obs) initiateMultipart(ctx context.Context, key string, options *multipartOptions) (uploadId string, err error) {
	var client *ObsClient
	var output *InitiateMultipartUploadOutput

	client, err = newObsClient(options.secret.Id, options.secret.Key, options.endpoint)
	if nil != err {
		return
	}

	objectOperationInput := &ObjectOperationInput{
		Bucket: options.bucket,
		Key:    key,
	}
	output, err = client.InitiateMultipartUpload(&InitiateMultipartUploadInput{
		ObjectOperationInput: *objectOperationInput,
	})

	if err != nil {
		return
	}
	uploadId = output.UploadId

	return
}

func (_ *_obs) completeMultipart(ctx context.Context, key string, uploadId string, objects []Object, options *multipartOptions) (err error) {
	var client *ObsClient

	client, err = newObsClient(options.secret.Id, options.secret.Key, options.endpoint)
	if nil != err {
		return
	}

	parts := make([]Part, 0, len(objects))
	for _, object := range objects {
		parts = append(parts, Part{
			ETag:       object.etag,
			PartNumber: int(object.size),
			Size:       object.size,
		})
	}

	input := &CompleteMultipartUploadInput{
		Bucket:   options.bucket,
		Key:      key,
		UploadId: uploadId,
		Parts:    parts,
	}
	_, err = client.CompleteMultipartUpload(input)

	return
}

func (_ *_obs) abortMultipart(ctx context.Context, key string, uploadId string, options *multipartOptions) (err error) {
	var client *ObsClient

	client, err = newObsClient(options.secret.Id, options.secret.Key, options.endpoint)
	if nil != err {
		return
	}

	input := &AbortMultipartUploadInput{
		Bucket:   options.bucket,
		Key:      key,
		UploadId: uploadId,
	}
	_, err = client.AbortMultipartUpload(input)

	return
}

func (_ *_obs) delete(ctx context.Context, key string, options *deleteOptions) (err error) {
	var client *ObsClient

	client, err = newObsClient(options.secret.Id, options.secret.Key, options.endpoint)
	if nil != err {
		return
	}

	input := &DeleteObjectInput{
		Bucket: options.bucket,
		Key:    key,
	}
	_, err = client.DeleteObject(input)

	return
}
