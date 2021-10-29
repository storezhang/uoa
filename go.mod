module github.com/storezhang/uoa

go 1.16

require (
	github.com/aws/aws-sdk-go v1.41.7
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/mozillazg/go-httpheader v0.3.0 // indirect
	github.com/storezhang/gox v1.7.2
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.209
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts v1.0.209
	github.com/tencentyun/cos-go-sdk-v5 v0.7.29
	golang.org/x/crypto v0.0.0-20210812204632-0ba0e8f03122 // indirect
	golang.org/x/text v0.3.7 // indirect
)

// replace github.com/storezhang/gox => ../gox
// replace github.com/storezhang/gox => ../pangu
