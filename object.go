package uoa

import (
	`github.com/tencentyun/cos-go-sdk-v5`
)

// Object 文件数据
type Object struct {
	key     string
	etag    string
	size    int64
	part    int32
	version string
}

// NewObject 创建文件数据
func NewObject(key string, etag string, size int64, part int32, version string) Object {
	return Object{
		key:     key,
		etag:    etag,
		size:    size,
		part:    part,
		version: version,
	}
}

func (o *Object) cos() cos.Object {
	return cos.Object{
		Key:        o.key,
		ETag:       o.etag,
		Size:       o.size,
		PartNumber: int(o.part),
		VersionId:  o.version,
	}
}
