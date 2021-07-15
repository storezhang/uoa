package uoa

import (
	`context`
	`net/url`
)

type uoaInternal interface {
	sts(ctx context.Context, options *stsOptions, keys ...string) (sts Sts, err error)
	url(ctx context.Context, key string, filename string, options *urlOptions) (downloadUrl *url.URL, err error)
}
