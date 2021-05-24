package uoa

import (
	`context`
	`net/url`
)

type uoaInternal interface {
	uploadUrl(ctx context.Context, key string, options *options) (uploadUrl *url.URL, err error)
	downloadUrl(ctx context.Context, key string, filename string, options *options) (downloadUrl *url.URL, err error)
}
