package uoa

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// 内部接口封装
// 使用模板方法设计模式
type template struct {
	cos executor
	s3  executor
}

func (t *template) Exist(ctx context.Context, bucket string, path Path, opts ...option) (exist bool, err error) {
	options := defaultOptions
	for _, opt := range opts {
		opt.apply(options)
	}

	key := t.key(path, options.environment, options.separator)
	switch options.uoaType {
	case TypeCos:
		exist, err = t.cos.exist(ctx, bucket, key, options)
	case TypeS3:
		exist, err = t.s3.exist(ctx, bucket, key, options)
	}

	return
}

func (t *template) Credentials(ctx context.Context, path Path, opts ...credentialsOption) (credentials *Credentials, err error) {
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
	case TypeS3:
		base, err = t.s3.credentials(ctx, options, key)
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

func (t *template) Url(ctx context.Context, bucket string, path Path, opts ...urlOption) (url *url.URL, err error) {
	options := defaultUrlOptions()
	for _, opt := range opts {
		opt.applyUrl(options)
	}

	key := t.key(path, options.environment, options.separator)
	switch options.uoaType {
	case TypeCos:
		url, err = t.cos.url(ctx, bucket, key, options)
	case TypeS3:
		url, err = t.s3.url(ctx, bucket, key, options)
	}

	return
}

func (t *template) InitiateMultipart(ctx context.Context, path Path, opts ...multipartOption) (uploadId string, err error) {
	options := defaultMultipartOptions()
	for _, opt := range opts {
		opt.applyMultipart(options)
	}

	fileKey := t.key(path, options.environment, options.separator)
	switch options.uoaType {
	case TypeCos:
		uploadId, err = t.cos.initiateMultipart(ctx, fileKey, options)
	case TypeS3:
		uploadId, err = t.s3.initiateMultipart(ctx, fileKey, options)
	}

	return
}

func (t *template) CompleteMultipart(ctx context.Context, path Path, uploadId string, objects []Object, opts ...multipartOption) (err error) {
	options := defaultMultipartOptions()
	for _, opt := range opts {
		opt.applyMultipart(options)
	}

	fileKey := t.key(path, options.environment, options.separator)
	switch options.uoaType {
	case TypeCos:
		err = t.cos.completeMultipart(ctx, fileKey, uploadId, objects, options)
	case TypeS3:
		err = t.s3.completeMultipart(ctx, fileKey, uploadId, objects, options)
	}

	return
}

func (t *template) AbortMultipart(ctx context.Context, path Path, uploadId string, opts ...multipartOption) (err error) {
	options := defaultMultipartOptions()
	for _, opt := range opts {
		opt.applyMultipart(options)
	}

	fileKey := t.key(path, options.environment, options.separator)
	switch options.uoaType {
	case TypeCos:
		err = t.cos.abortMultipart(ctx, fileKey, uploadId, options)
	case TypeS3:
		err = t.s3.abortMultipart(ctx, fileKey, uploadId, options)
	}

	return
}

func (t *template) Delete(ctx context.Context, path Path, opts ...deleteOption) (err error) {
	options := defaultDeleteOptions()
	for _, opt := range opts {
		opt.applyDelete(options)
	}

	key := t.key(path, options.environment, options.separator)
	switch options.uoaType {
	case TypeCos:
		err = t.cos.delete(ctx, key, options)
	case TypeS3:
		err = t.s3.delete(ctx, key, options)
	}

	return
}

func (t *template) key(path Path, environment string, separator string) (key string) {
	paths := path.Paths()
	if "" != environment {
		paths = append([]string{environment}, paths...)
	}
	key = strings.Join(paths, separator)

	return
}
