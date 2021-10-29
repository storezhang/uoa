package uoa

import (
	`strings`
	`sync/atomic`
)

var emptySecurityHolder = securityHolder{}

type securityHolder struct {
	accessKey     string
	securityKey   string
	securityToken string
}

type securityProvider interface {
	getSecurity() securityHolder
}

type BasicSecurityProvider struct {
	val atomic.Value
}

func (bsp *BasicSecurityProvider) getSecurity() securityHolder {
	if sh, ok := bsp.val.Load().(securityHolder); ok {
		return sh
	}
	return emptySecurityHolder
}

func (bsp *BasicSecurityProvider) refresh(accessKey, securityKey, securityToken string) {
	bsp.val.Store(securityHolder{
		accessKey:     strings.TrimSpace(accessKey),
		securityKey:   strings.TrimSpace(securityKey),
		securityToken: strings.TrimSpace(securityToken),
	})
}

func NewBasicSecurityProvider(accessKey, securityKey, securityToken string) *BasicSecurityProvider {
	bsp := &BasicSecurityProvider{}
	bsp.refresh(accessKey, securityKey, securityToken)
	return bsp
}
