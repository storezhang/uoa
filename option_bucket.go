package uoa

var _ deleteOption = (*optionVersion)(nil)

type optionBucket struct {
	bucket string
}

// BucketName 桶名称
func Bucket(bucket string) *optionBucket {
	return &optionBucket{
		bucket: bucket,
	}
}

func (v *optionBucket) applyDelete(options *deleteOptions) {
	options.bucket = v.bucket
}

func (v *optionBucket) applyMultipart(options *multipartOptions) {
	options.bucket = v.bucket
}
