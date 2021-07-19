package uoa

import (
	`github.com/storezhang/gox`
)

var _ urlOption = (*optionEnvironment)(nil)

type optionEnvironment struct {
	environment string
}

// Environment 配置应用名称
func Environment(environment gox.Environment) *optionEnvironment {
	return &optionEnvironment{
		environment: string(environment),
	}
}

func (e *optionEnvironment) apply(options *options) {
	options.environment = e.environment
}

func (e *optionEnvironment) applyUrl(options *urlOptions) {
	options.environment = e.environment
}

func (e *optionEnvironment) applyCredential(options *credentialsOptions) {
	options.environment = e.environment
}
