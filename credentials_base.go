package uoa

import (
	`time`
)

type credentialsBase struct {
	// 临时授权，相当于用户名
	Id string `json:"id" yaml:"id" xml:"id"`
	// 临时授权，相当于密码
	Key string `json:"key" yaml:"key" xml:"key"`
	// 临时授权
	Token string `json:"token" yaml:"token" xml:"token"`
	// 过期时间
	Expired time.Time `json:"expired" yaml:"expired" xml:"expired"`
}
