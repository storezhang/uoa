package uoa

import (
	`time`
)

var _ urlOption = (*optionExpired)(nil)

type optionExpired struct {
	expired time.Duration
}

// Expired 配置应用名称
func Expired(expired time.Duration) *optionExpired {
	return &optionExpired{
		expired: expired,
	}
}

func (e *optionExpired) apply(options *options) {
	options.expired = e.expired
}

func (e *optionExpired) applyUrl(options *urlOptions) {
	options.expired = e.expired
}
