package uoa

import (
	`github.com/storezhang/gox`
)

var _ option = (*optionSecret)(nil)

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

func (b *optionSecret) apply(options *options) {
	options.secret.Id = b.id
	options.secret.Key = b.key
}
