package uoa

import (
	`github.com/storezhang/pangu`
)

func init() {
	app := pangu.New()

	if err := app.Provides(
		NewCos,
		New,
	); nil != err {
		panic(err)
	}
}
