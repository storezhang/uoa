package uoa

import (
	`encoding/xml`
	`io`
	`net/http`
	`time`
)

type ISseHeader interface {
	GetEncryption() string
	GetKey() string
}

type IBaseModel interface {
	setStatusCode(statusCode int)

	setRequestID(requestID string)

	setResponseHeaders(responseHeaders map[string][]string)
}

type ISerializable interface {
	trans(isObs bool) (map[string]string, map[string][]string, interface{}, error)
}

type IReadCloser interface {
	setReadCloser(body io.ReadCloser)
}

type BaseModel struct {
	StatusCode      int                 `xml:"-"`
	RequestId       string              `xml:"RequestId" json:"request_id"`
	ResponseHeaders map[string][]string `xml:"-"`
}

type GetObjectOutput struct {
	GetObjectMetadataOutput
	DeleteMarker       bool
	CacheControl       string
	ContentDisposition string
	ContentEncoding    string
	ContentLanguage    string
	Expires            string
	Body               io.ReadCloser
}

type GetObjectMetadataOutput struct {
	BaseModel
	VersionId               string
	WebsiteRedirectLocation string
	Expiration              string
	Restore                 string
	ObjectType              string
	NextAppendPosition      string
	StorageClass            StorageClassType
	ContentLength           int64
	ContentType             string
	ETag                    string
	AllowOrigin             string
	AllowHeader             string
	AllowMethod             string
	ExposeHeader            string
	MaxAgeSeconds           int
	LastModified            time.Time
	SseHeader               ISseHeader
	Metadata                map[string]string
}

func (baseModel *BaseModel) setStatusCode(statusCode int) {
	baseModel.StatusCode = statusCode
}

func (baseModel *BaseModel) setRequestID(requestID string) {
	baseModel.RequestId = requestID
}

func (baseModel *BaseModel) setResponseHeaders(responseHeaders map[string][]string) {
	baseModel.ResponseHeaders = responseHeaders
}

type SseCHeader struct {
	Encryption string
	Key        string
	KeyMD5     string
}

func (header SseCHeader) GetEncryption() string {
	if header.Encryption != "" {
		return header.Encryption
	}
	return "AES256"
}

func (header SseCHeader) GetKey() string {
	return header.Key
}

type CreateSignedUrlInput struct {
	Method      HttpMethodType
	Bucket      string
	Key         string
	SubResource SubResourceType
	Expires     int
	Headers     map[string]string
	QueryParams map[string]string
}

type CreateSignedUrlOutput struct {
	SignedUrl                  string
	ActualSignedRequestHeaders http.Header
}

type ObjectOperationInput struct {
	Bucket                  string
	Key                     string
	ACL                     AclType
	GrantReadId             string
	GrantReadAcpId          string
	GrantWriteAcpId         string
	GrantFullControlId      string
	StorageClass            StorageClassType
	WebsiteRedirectLocation string
	Expires                 int64
	SseHeader               ISseHeader
	Metadata                map[string]string
}

type HeadObjectInput struct {
	Bucket    string
	Key       string
	VersionId string
}

type InitiateMultipartUploadInput struct {
	ObjectOperationInput
	ContentType  string
	EncodingType string
}

type InitiateMultipartUploadOutput struct {
	BaseModel
	XMLName      xml.Name `xml:"InitiateMultipartUploadResult"`
	Bucket       string   `xml:"Bucket"`
	Key          string   `xml:"Key"`
	UploadId     string   `xml:"UploadId"`
	SseHeader    ISseHeader
	EncodingType string `xml:"EncodingType,omitempty"`
}

type CompleteMultipartUploadInput struct {
	Bucket       string   `xml:"-"`
	Key          string   `xml:"-"`
	UploadId     string   `xml:"-"`
	XMLName      xml.Name `xml:"CompleteMultipartUpload"`
	Parts        []Part   `xml:"Part"`
	EncodingType string   `xml:"-"`
}

type CompleteMultipartUploadOutput struct {
	BaseModel
	VersionId    string     `xml:"-"`
	SseHeader    ISseHeader `xml:"-"`
	XMLName      xml.Name   `xml:"CompleteMultipartUploadResult"`
	Location     string     `xml:"Location"`
	Bucket       string     `xml:"Bucket"`
	Key          string     `xml:"Key"`
	ETag         string     `xml:"ETag"`
	EncodingType string     `xml:"EncodingType,omitempty"`
}

type AbortMultipartUploadInput struct {
	Bucket   string
	Key      string
	UploadId string
}

type DeleteObjectInput struct {
	Bucket    string
	Key       string
	VersionId string
}

type DeleteObjectOutput struct {
	BaseModel
	VersionId    string
	DeleteMarker bool
}

type Part struct {
	XMLName      xml.Name  `xml:"Part"`
	PartNumber   int       `xml:"PartNumber"`
	ETag         string    `xml:"ETag"`
	LastModified time.Time `xml:"LastModified,omitempty"`
	Size         int64     `xml:"Size,omitempty"`
}

type partSlice []Part

func (parts partSlice) Len() int {
	return len(parts)
}

func (parts partSlice) Less(i, j int) bool {
	return parts[i].PartNumber < parts[j].PartNumber
}

func (parts partSlice) Swap(i, j int) {
	parts[i], parts[j] = parts[j], parts[i]
}

type SseKmsHeader struct {
	Encryption string
	Key        string
	isObs      bool
}

func (header SseKmsHeader) GetEncryption() string {
	if header.Encryption != "" {
		return header.Encryption
	}
	if !header.isObs {
		return "aws:kms"
	}
	return "kms"
}

func (header SseKmsHeader) GetKey() string {
	return header.Key
}
