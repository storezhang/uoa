package uoa

import (
	`github.com/storezhang/gox`
)

var _ option = (*optionEnvironment)(nil)

type optionEnvironment struct {
	environment string
}

// Environment 配置应用名称
func Environment(environment gox.Environment) *optionEnvironment {
	return &optionEnvironment{
		environment: string(environment),
	}
}

func (b *optionEnvironment) apply(options *options) {
	options.environment = b.environment
}
