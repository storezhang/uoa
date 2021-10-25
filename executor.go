package uoa

import (
	`context`
	`net/url`
)

type executor interface {
	exist(ctx context.Context, bucket string, key string, options *options) (exist bool, err error)

	credentials(ctx context.Context, options *credentialsOptions, keys ...string) (credentials *credentialsBase, err error)

	url(ctx context.Context, bucket string, key string, options *urlOptions) (downloadUrl *url.URL, err error)

	initiateMultipart(ctx context.Context, key string, options *multipartOptions) (uploadId string, err error)

	completeMultipart(ctx context.Context, key string, uploadId string, objects []Object, options *multipartOptions) (err error)

	abortMultipart(ctx context.Context, key string, uploadId string, options *multipartOptions) (err error)

	delete(ctx context.Context, key string, options *deleteOptions) (err error)
}
