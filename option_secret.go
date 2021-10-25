package uoa

import (
	`github.com/storezhang/gox`
)

var _ urlOption = (*optionSecret)(nil)

type optionSecret struct {
	id      string
	key     string
	uoaType Type
}

// Secret 配置授权
func Secret(secret gox.Secret) *optionSecret {
	return &optionSecret{
		id:  secret.Id,
		key: secret.Key,
	}
}

// Tencentyun 配置腾讯云授权
func Tencentyun(secretId string, secretKey string) *optionSecret {
	return &optionSecret{
		id:      secretId,
		key:     secretKey,
		uoaType: TypeCos,
	}
}

// S3 配置授权
func S3(secretId string, secretKey string) *optionSecret {
	return &optionSecret{
		id:      secretId,
		key:     secretKey,
		uoaType: TypeS3,
	}
}

func (s *optionSecret) apply(options *options) {
	options.secret.Id = s.id
	options.secret.Key = s.key
	options.uoaType = s.uoaType
}

func (s *optionSecret) applyUrl(options *urlOptions) {
	options.secret.Id = s.id
	options.secret.Key = s.key
	options.uoaType = s.uoaType
}

func (s *optionSecret) applyCredential(options *credentialsOptions) {
	options.secret.Id = s.id
	options.secret.Key = s.key
	options.uoaType = s.uoaType
}
