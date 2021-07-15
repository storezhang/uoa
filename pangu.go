package uoa

import (
	`github.com/storezhang/pangu`
)

func init() {
	if err := pangu.New().Provides(
		NewCos,
		New,
	); nil != err {
		panic(err)
	}
}
