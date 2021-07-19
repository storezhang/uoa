package uoa

import (
	`context`
	`fmt`
	`net/url`
	`strings`
)

// 内部接口封装
// 使用模板方法设计模式
type uoaTemplate struct {
	cos uoaInternal
}

func (t *uoaTemplate) Credentials(ctx context.Context, path Path, opts ...credentialsOption) (credentials *Credentials, err error) {
	options := defaultCredentialOptions()
	for _, opt := range opts {
		opt.applyCredential(options)
	}

	key := t.key(path, options.environment, options.separator)
	var keys []string
	if 0 != len(options.patterns) {
		keys = make([]string, 0, len(options.patterns))
		for _, pattern := range options.patterns {
			keys = append(keys, fmt.Sprintf("%s%s%s", key, options.separator, pattern))
		}
	} else {
		keys = []string{key}
	}

	var base *credentialsBase
	switch options.uoaType {
	case TypeCos:
		base, err = t.cos.credentials(ctx, options, keys...)
	}
	if nil != err {
		return
	}

	// 注入通用字段
	credentials = &Credentials{
		credentialsBase: base,
		Url:             options.endpoint,
		Separator:       options.separator,
	}

	return
}

func (t *uoaTemplate) Url(ctx context.Context, path Path, filename string, opts ...urlOption) (url *url.URL, err error) {
	options := defaultUrlOptions()
	for _, opt := range opts {
		opt.applyUrl(options)
	}

	key := t.key(path, options.environment, options.separator)
	switch options.uoaType {
	case TypeCos:
		url, err = t.cos.url(ctx, key, filename, options)
	}

	return
}

func (t *uoaTemplate) key(path Path, environment string, separator string) (key string) {
	paths := path.Paths()
	if "" != environment {
		paths = append([]string{environment}, paths...)
	}
	key = strings.Join(path.Paths(), separator)

	return
}
