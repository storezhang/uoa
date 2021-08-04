package uoa

import (
	`context`
	`net/url`

	`github.com/tencentyun/cos-go-sdk-v5`
)

type uoaInternal interface {
	credentials(ctx context.Context, options *credentialsOptions, keys ...string) (credentials *credentialsBase, err error)
	url(ctx context.Context, key string, filename string, options *urlOptions) (downloadUrl *url.URL, err error)
	initiateMultipartUpload(ctx context.Context, key string, options *multipartOptions) (uploadId string, err error)
	completeMultipartUpload(ctx context.Context, key string, uploadId string, parts []cos.Object, options *multipartOptions) (err error)
	abortMultipartUpload(ctx context.Context, key string, uploadId string, options *multipartOptions) (err error)
	delete(ctx context.Context, key string, options *deleteOptions) (err error)
}
