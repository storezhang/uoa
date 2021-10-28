package uoa

import (
	"errors"
	"net/http"
	"sort"
	"strings"
)

type ObsClient struct {
	conf       *config
	httpClient *http.Client
}

func NewObsClient(accessKey, securityKey, endPoint string, configures ...configuror) (client *ObsClient, err error) {
	conf := &config{endPoint: endPoint}
	conf.securityProviders = make([]securityProvider, 0, 3)
	conf.securityProviders = append(conf.securityProviders, NewBasicSecurityProvider(accessKey, securityKey, ""))

	conf.maxRetryCount = -1
	conf.maxRedirectCount = -1
	for _, configure := range configures {
		configure(conf)
	}

	if err = conf.initConfigWithDefault(); err != nil {
		return
	}

	if err = conf.getTransport(); nil != err {
		return
	}

	client = &ObsClient{
		conf: conf,
		httpClient: &http.Client{
			Transport:     conf.transport,
			CheckRedirect: checkRedirectFunc,
		},
	}

	return
}

func (o ObsClient) CreateSignedUrl(input *CreateSignedUrlInput) (output *CreateSignedUrlOutput, err error) {
	if input == nil {
		return nil, errors.New("CreateSignedUrlInput is nil")
	}

	params := make(map[string]string, len(input.QueryParams))
	for key, value := range input.QueryParams {
		params[key] = value
	}

	if input.SubResource != "" {
		params[string(input.SubResource)] = ""
	}

	headers := make(map[string][]string, len(input.Headers))
	for key, value := range input.Headers {
		headers[key] = []string{value}
	}

	if input.Expires <= 0 {
		input.Expires = 300
	}

	requestURL, err := o.doAuthTemporary(string(input.Method), input.Bucket, input.Key, params, headers, int64(input.Expires))
	if err != nil {
		return nil, err
	}

	output = &CreateSignedUrlOutput{
		SignedUrl:                  requestURL,
		ActualSignedRequestHeaders: headers,
	}

	return
}

func (o ObsClient) getSecurity() securityHolder {
	if o.conf.securityProviders != nil {
		for _, sp := range o.conf.securityProviders {
			if sp == nil {
				continue
			}
			sh := sp.getSecurity()
			if sh.accessKey != "" && sh.securityKey != "" {
				return sh
			}
		}
	}

	return emptySecurityHolder
}

func (o ObsClient) Close() {
	o.httpClient = nil
	o.conf.transport.CloseIdleConnections()
	o.conf = nil
}

func (o ObsClient) HeadBucket(bucket string, extensions ...extensionOptions) (resp *http.Response, err error) {
	var redirectURL string
	var requestURL string
	var req *http.Request

	params := make(map[string]string)
	headers := make(map[string][]string)

	for _, extension := range extensions {
		if extensionHeader, ok := extension.(extensionHeaders); ok {
			_err := extensionHeader(headers, true)
			if _err != nil {

			}
		} else {

		}
	}

	_data, _err := prepareData(headers, nil)
	if _err != nil {
		return nil, _err
	}

	req, err = o.getRequest(redirectURL, requestURL, false, _data, "HEAD", bucket, "", params, headers)
	if nil == req {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err = o.httpClient.Do(req)

	return
}

func (o ObsClient) HeadObject(input *HeadObjectInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	var redirectURL string
	var requestURL string
	var req *http.Request
	var resp *http.Response
	var data interface{}

	params := make(map[string]string)
	headers := make(map[string][]string)

	for _, extension := range extensions {
		if extensionHeader, ok := extension.(extensionHeaders); ok {
			_err := extensionHeader(headers, true)
			if _err != nil {

			}
		} else {

		}
	}
	//准备参数、请求头、数据
	params, headers, data, err = input.trans(true)
	headers = prepareHeaders(headers, false, true)

	_data, _err := prepareData(headers, data)
	if _err != nil {
		return nil, _err
	}

	// 参数校验
	if input == nil {
		return nil, errors.New("InitiateMultipartUploadInput is nil")
	}
	if strings.TrimSpace(input.Bucket) == "" && !o.conf.cname {
		err = errors.New("Bucket is empty")
		return
	}
	if strings.TrimSpace(input.Key) == "" {
		err = errors.New("Key is empty")
		return
	}

	// 构造HttpRequest
	req, err = o.getRequest(redirectURL, requestURL, false, _data, "POST", input.Bucket, input.Key, params, headers)
	if nil == req {
		return
	}

	// 发送POST请求
	resp, err = o.httpClient.Do(req)
	// 解析Response到output中
	err = ParseResponseToBaseModel(resp, output, true, true)
	if nil != err {
		output = nil
	}

	return
}

func (o ObsClient) InitiateMultipartUpload(input *InitiateMultipartUploadInput, extensions ...extensionOptions) (output *InitiateMultipartUploadOutput, err error) {
	var redirectURL string
	var requestURL string
	var req *http.Request
	var resp *http.Response
	var data interface{}

	params := make(map[string]string)
	headers := make(map[string][]string)

	for _, extension := range extensions {
		if extensionHeader, ok := extension.(extensionHeaders); ok {
			_err := extensionHeader(headers, true)
			if _err != nil {

			}
		} else {

		}
	}
	//准备参数、请求头、数据
	params, headers, data, err = input.trans(true)
	headers = prepareHeaders(headers, false, true)

	_data, _err := prepareData(headers, data)
	if _err != nil {
		return nil, _err
	}

	// 参数校验
	if input == nil {
		return nil, errors.New("InitiateMultipartUploadInput is nil")
	}
	if strings.TrimSpace(input.Bucket) == "" && !o.conf.cname {
		err = errors.New("Bucket is empty")
		return
	}
	if strings.TrimSpace(input.Key) == "" {
		err = errors.New("Key is empty")
		return
	}
	if input.ContentType == "" && input.Key != "" {
		if contentType, ok := mimeTypes[strings.ToLower(input.Key[strings.LastIndex(input.Key, ".")+1:])]; ok {
			input.ContentType = contentType
		}
	}

	// 构造HttpRequest
	req, err = o.getRequest(redirectURL, requestURL, false, _data, "POST", input.Bucket, input.Key, params, headers)
	if nil == req {
		return
	}

	// 发送POST请求
	resp, err = o.httpClient.Do(req)
	// 解析Response到output中
	err = ParseResponseToBaseModel(resp, output, true, true)
	if nil != err {
		output = nil
	}
	ParseInitiateMultipartUploadOutput(output)
	if output.EncodingType == "url" {
		err = decodeInitiateMultipartUploadOutput(output)
		if err != nil {
			output = nil
		}
	}

	return
}

func (o ObsClient) CompleteMultipartUpload(input *CompleteMultipartUploadInput, extensions ...extensionOptions) (output *CompleteMultipartUploadOutput, err error) {
	var redirectURL string
	var requestURL string
	var req *http.Request
	var resp *http.Response
	var data interface{}

	params := make(map[string]string)
	headers := make(map[string][]string)

	for _, extension := range extensions {
		if extensionHeader, ok := extension.(extensionHeaders); ok {
			_err := extensionHeader(headers, true)
			if _err != nil {

			}
		} else {

		}
	}

	// 参数校验
	if input == nil {
		return nil, errors.New("CompleteMultipartUploadInput is nil")
	}
	if input.UploadId == "" {
		return nil, errors.New("UploadId is empty")
	}
	if strings.TrimSpace(input.Bucket) == "" && !o.conf.cname {
		err = errors.New("Bucket is empty")
		return
	}
	if strings.TrimSpace(input.Key) == "" {
		err = errors.New("Key is empty")
		return
	}

	var parts partSlice = input.Parts
	sort.Sort(parts)
	//准备参数、请求头、数据
	params, headers, data, err = input.trans(true)
	headers = prepareHeaders(headers, false, true)

	_data, _err := prepareData(headers, data)
	if _err != nil {
		return nil, _err
	}

	// 构造HttpRequest
	req, err = o.getRequest(redirectURL, requestURL, false, _data, "POST", input.Bucket, input.Key, params, headers)

	// 发送POST请求
	resp, err = o.httpClient.Do(req)
	// 解析Response到output中
	err = ParseResponseToBaseModel(resp, output, true, true)
	if nil != err {
		output = nil
		return
	}
	ParseCompleteMultipartUploadOutput(output)
	if output.EncodingType == "url" {
		err = decodeCompleteMultipartUploadOutput(output)
		if err != nil {
			output = nil
		}
	}

	return
}

func (o ObsClient) AbortMultipartUpload(input *AbortMultipartUploadInput, extensions ...extensionOptions) (output *BaseModel, err error) {
	var redirectURL string
	var requestURL string
	var req *http.Request
	var resp *http.Response
	var data interface{}

	params := make(map[string]string)
	headers := make(map[string][]string)

	for _, extension := range extensions {
		if extensionHeader, ok := extension.(extensionHeaders); ok {
			_err := extensionHeader(headers, true)
			if _err != nil {

			}
		} else {

		}
	}
	// 参数校验
	if input == nil {
		return nil, errors.New("AbortMultipartUploadInput is nil")
	}
	if input.UploadId == "" {
		return nil, errors.New("UploadId is empty")
	}
	if strings.TrimSpace(input.Bucket) == "" && !o.conf.cname {
		err = errors.New("Bucket is empty")
		return
	}
	if strings.TrimSpace(input.Key) == "" {
		err = errors.New("Key is empty")
		return
	}

	//准备参数、请求头、数据
	params, headers, data, err = input.trans(true)
	headers = prepareHeaders(headers, false, true)

	_data, _err := prepareData(headers, data)
	if _err != nil {
		return nil, _err
	}

	req, err = o.getRequest(redirectURL, requestURL, false, _data, "DELETE", input.Bucket, input.Key, params, headers)

	// 发送POST请求
	resp, err = o.httpClient.Do(req)
	// 解析Response到output中
	err = ParseResponseToBaseModel(resp, output, true, true)

	return
}

func (o ObsClient) DeleteObject(input *DeleteObjectInput, extensions ...extensionOptions) (output *DeleteObjectOutput, err error) {
	var redirectURL string
	var requestURL string
	var req *http.Request
	var resp *http.Response
	var data interface{}

	params := make(map[string]string)
	headers := make(map[string][]string)

	for _, extension := range extensions {
		if extensionHeader, ok := extension.(extensionHeaders); ok {
			_err := extensionHeader(headers, true)
			if _err != nil {

			}
		} else {

		}
	}
	// 参数校验
	if input == nil {
		return nil, errors.New("AbortMultipartUploadInput is nil")
	}
	if strings.TrimSpace(input.Bucket) == "" && !o.conf.cname {
		err = errors.New("Bucket is empty")
		return
	}
	if strings.TrimSpace(input.Key) == "" {
		err = errors.New("Key is empty")
		return
	}

	//准备参数、请求头、数据
	params, headers, data, err = input.trans(true)
	headers = prepareHeaders(headers, false, true)

	_data, _err := prepareData(headers, data)
	if _err != nil {
		return nil, _err
	}

	req, err = o.getRequest(redirectURL, requestURL, false, _data, "DELETE", input.Bucket, input.Key, params, headers)

	// 发送POST请求
	resp, err = o.httpClient.Do(req)
	// 解析Response到output中
	err = ParseResponseToBaseModel(resp, output, true, true)
	if nil != err {
		output = nil
		return
	}
	ParseDeleteObjectOutput(output)

	return
}
