package uoa

import (
	`encoding/xml`
	`io`
	`net/http`
	`time`
)

type isseHeader interface {
	getEncryption() string
	getKey() string
}

type ibaseModel interface {
	setStatusCode(statusCode int)

	setRequestID(requestID string)

	setResponseHeaders(responseHeaders map[string][]string)
}

type iserializable interface {
	trans(isObs bool) (map[string]string, map[string][]string, interface{}, error)
}

type ireadCloser interface {
	setReadCloser(body io.ReadCloser)
}

type baseModel struct {
	statusCode      int                 `xml:"-"`
	requestId       string              `xml:"RequestId" json:"request_id"`
	responseHeaders map[string][]string `xml:"-"`
}

type getObjectOutput struct {
	getObjectMetadataOutput
	deleteMarker       bool
	cacheControl       string
	contentDisposition string
	contentEncoding    string
	contentLanguage    string
	expires            string
	body               io.ReadCloser
}

type getObjectMetadataOutput struct {
	baseModel
	versionId               string
	websiteRedirectLocation string
	expiration              string
	restore                 string
	objectType              string
	nextAppendPosition      string
	storageClass            StorageClassType
	contentLength           int64
	contentType             string
	eTag                    string
	allowOrigin             string
	allowHeader             string
	allowMethod             string
	exposeHeader            string
	maxAgeSeconds           int
	lastModified            time.Time
	sseHeader               isseHeader
	metadata                map[string]string
}

func (baseModel *baseModel) setStatusCode(statusCode int) {
	baseModel.statusCode = statusCode
}

func (baseModel *baseModel) setRequestID(requestID string) {
	baseModel.requestId = requestID
}

func (baseModel *baseModel) setResponseHeaders(responseHeaders map[string][]string) {
	baseModel.responseHeaders = responseHeaders
}

type sseCHeader struct {
	encryption string
	key        string
	keyMD5     string
}

func (header sseCHeader) getEncryption() string {
	if header.encryption != "" {
		return header.encryption
	}
	return "AES256"
}

func (header sseCHeader) getKey() string {
	return header.key
}

type createSignedUrlInput struct {
	method      HttpMethodType
	bucket      string
	key         string
	subResource SubResourceType
	expires     int
	headers     map[string]string
	queryParams map[string]string
}

type createSignedUrlOutput struct {
	signedUrl                  string
	actualSignedRequestHeaders http.Header
}

type objectOperationInput struct {
	bucket                  string
	key                     string
	acl                     AclType
	grantReadId             string
	grantReadAcpId          string
	grantWriteAcpId         string
	grantFullControlId      string
	storageClass            StorageClassType
	websiteRedirectLocation string
	expires                 int64
	sseHeader               isseHeader
	metadata                map[string]string
}

type headObjectInput struct {
	bucket    string
	key       string
	versionId string
}

type initiateMultipartUploadInput struct {
	objectOperationInput
	contentType  string
	encodingType string
}

type initiateMultipartUploadOutput struct {
	baseModel
	xmlName      xml.Name `xml:"InitiateMultipartUploadResult"`
	bucket       string   `xml:"Bucket"`
	key          string   `xml:"Key"`
	uploadId     string   `xml:"UploadId"`
	sseHeader    isseHeader
	encodingType string `xml:"EncodingType,omitempty"`
}

type completeMultipartUploadInput struct {
	bucket       string   `xml:"-"`
	key          string   `xml:"-"`
	uploadId     string   `xml:"-"`
	xmlName      xml.Name `xml:"CompleteMultipartUpload"`
	parts        []part   `xml:"Part"`
	encodingType string   `xml:"-"`
}

type completeMultipartUploadOutput struct {
	baseModel
	versionId    string     `xml:"-"`
	sseHeader    isseHeader `xml:"-"`
	xmlName      xml.Name   `xml:"CompleteMultipartUploadResult"`
	location     string     `xml:"Location"`
	bucket       string     `xml:"Bucket"`
	key          string     `xml:"Key"`
	eTag         string     `xml:"ETag"`
	encodingType string     `xml:"EncodingType,omitempty"`
}

type abortMultipartUploadInput struct {
	bucket   string
	key      string
	uploadId string
}

type deleteObjectInput struct {
	bucket    string
	key       string
	versionId string
}

type deleteObjectOutput struct {
	baseModel
	versionId    string
	deleteMarker bool
}

type part struct {
	xmlName      xml.Name  `xml:"Part"`
	partNumber   int       `xml:"PartNumber"`
	eTag         string    `xml:"ETag"`
	lastModified time.Time `xml:"LastModified,omitempty"`
	size         int64     `xml:"Size,omitempty"`
}

type partSlice []part

func (parts partSlice) Len() int {
	return len(parts)
}

func (parts partSlice) Less(i, j int) bool {
	return parts[i].partNumber < parts[j].partNumber
}

func (parts partSlice) Swap(i, j int) {
	parts[i], parts[j] = parts[j], parts[i]
}

type sseKmsHeader struct {
	encryption string
	key        string
	isObs      bool
}

func (header sseKmsHeader) getEncryption() string {
	if header.encryption != "" {
		return header.encryption
	}
	if !header.isObs {
		return "aws:kms"
	}
	return "kms"
}

func (header sseKmsHeader) getKey() string {
	return header.key
}
