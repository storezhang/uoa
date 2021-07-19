package uoa

import (
	`context`
	`net/url`
)

type uoaInternal interface {
	credentials(ctx context.Context, options *credentialsOptions, keys ...string) (credentials *credentialsBase, err error)
	url(ctx context.Context, key string, filename string, options *urlOptions) (downloadUrl *url.URL, err error)
}
