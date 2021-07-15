package uoa

import (
	`context`
	`net/url`
)

type uoaInternal interface {
	sts(ctx context.Context, key string, options *stsOptions) (sts Sts, err error)
	downloadUrl(ctx context.Context, key string, filename string, options *urlOptions) (downloadUrl *url.URL, err error)
}
