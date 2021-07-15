package uoa

import (
	`github.com/storezhang/gox`
)

var _ urlOption = (*optionSecret)(nil)

type optionSecret struct {
	// 授权，类似于用户名
	id string
	// 授权，类似于密码
	key string
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
	return Secret(gox.Secret{
		Id:  secretId,
		Key: secretKey,
	})
}

func (s *optionSecret) applyUrl(options *urlOptions) {
	options.secret.Id = s.id
	options.secret.Key = s.key
}

func (s *optionSecret) applySts(options *stsOptions) {
	options.secret.Id = s.id
	options.secret.Key = s.key
}
