package uoa

import (
	`errors`
	`net/http`
	`sort`
	`strings`
)

type obsClient struct {
	conf       *config
	httpClient *http.Client
}

func NewObsClient(accessKey, securityKey, endPoint string, configures ...configuror) (client *obsClient, err error) {
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

	client = &obsClient{
		conf: conf,
		httpClient: &http.Client{
			Transport:     conf.transport,
			CheckRedirect: checkRedirectFunc,
		},
	}

	return
}

func (o obsClient) createSignedUrl(input *createSignedUrlInput) (output *createSignedUrlOutput, err error) {
	if input == nil {
		return nil, errors.New("createSignedUrlInput is nil")
	}

	params := make(map[string]string, len(input.queryParams))
	for key, value := range input.queryParams {
		params[key] = value
	}

	if input.subResource != "" {
		params[string(input.subResource)] = ""
	}

	headers := make(map[string][]string, len(input.headers))
	for key, value := range input.headers {
		headers[key] = []string{value}
	}

	if input.expires <= 0 {
		input.expires = 300
	}

	requestURL, err := o.doAuthTemporary(string(input.method), input.bucket, input.key, params, headers, int64(input.expires))
	if err != nil {
		return nil, err
	}

	output = &createSignedUrlOutput{
		signedUrl:                  requestURL,
		actualSignedRequestHeaders: headers,
	}

	return
}

func (o obsClient) getSecurity() securityHolder {
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

func (o obsClient) Close() {
	o.httpClient = nil
	o.conf.transport.CloseIdleConnections()
	o.conf = nil
}

func (o obsClient) headBucket(bucket string, extensions ...extensionOptions) (resp *http.Response, err error) {
	var (
		redirectURL string
		requestURL  string
		req         *http.Request
	)

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

func (o obsClient) headObject(input *headObjectInput, extensions ...extensionOptions) (output *baseModel, err error) {
	var (
		redirectURL string
		requestURL  string
		req         *http.Request
		resp        *http.Response
		data        interface{}
	)

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
		return nil, errors.New("initiateMultipartUploadInput is nil")
	}
	if strings.TrimSpace(input.bucket) == "" && !o.conf.cname {
		err = errors.New("bucket is empty")
		return
	}
	if strings.TrimSpace(input.key) == "" {
		err = errors.New("key is empty")
		return
	}

	// 构造HttpRequest
	req, err = o.getRequest(redirectURL, requestURL, false, _data, "POST", input.bucket, input.key, params, headers)
	if nil == req {
		return
	}

	// 发送POST请求
	resp, err = o.httpClient.Do(req)
	// 解析Response到output中
	err = parseResponseToBaseModel(resp, output, true, true)
	if nil != err {
		output = nil
	}

	return
}

func (o obsClient) initiateMultipartUpload(input *initiateMultipartUploadInput, extensions ...extensionOptions) (output *initiateMultipartUploadOutput, err error) {
	var (
		redirectURL string
		requestURL  string
		req         *http.Request
		resp        *http.Response
		data        interface{}
	)

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
		return nil, errors.New("initiateMultipartUploadInput is nil")
	}
	if strings.TrimSpace(input.bucket) == "" && !o.conf.cname {
		err = errors.New("bucket is empty")
		return
	}
	if strings.TrimSpace(input.key) == "" {
		err = errors.New("key is empty")
		return
	}
	if input.contentType == "" && input.key != "" {
		if contentType, ok := mimeTypes[strings.ToLower(input.key[strings.LastIndex(input.key, ".")+1:])]; ok {
			input.contentType = contentType
		}
	}

	// 构造HttpRequest
	req, err = o.getRequest(redirectURL, requestURL, false, _data, "POST", input.bucket, input.key, params, headers)
	if nil == req {
		return
	}

	// 发送POST请求
	resp, err = o.httpClient.Do(req)
	// 解析Response到output中
	err = parseResponseToBaseModel(resp, output, true, true)
	if nil != err {
		output = nil
	}
	parseInitiateMultipartUploadOutput(output)
	if output.encodingType == "url" {
		err = decodeInitiateMultipartUploadOutput(output)
		if err != nil {
			output = nil
		}
	}

	return
}

func (o obsClient) completeMultipartUpload(input *completeMultipartUploadInput, extensions ...extensionOptions) (output *completeMultipartUploadOutput, err error) {
	var (
		redirectURL string
		requestURL  string
		req         *http.Request
		resp        *http.Response
		data        interface{}
	)

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
		return nil, errors.New("completeMultipartUploadInput is nil")
	}
	if input.uploadId == "" {
		return nil, errors.New("uploadId is empty")
	}
	if strings.TrimSpace(input.bucket) == "" && !o.conf.cname {
		err = errors.New("bucket is empty")
		return
	}
	if strings.TrimSpace(input.key) == "" {
		err = errors.New("key is empty")
		return
	}

	var parts partSlice = input.parts
	sort.Sort(parts)
	//准备参数、请求头、数据
	params, headers, data, err = input.trans(true)
	headers = prepareHeaders(headers, false, true)

	_data, _err := prepareData(headers, data)
	if _err != nil {
		return nil, _err
	}

	// 构造HttpRequest
	req, err = o.getRequest(redirectURL, requestURL, false, _data, "POST", input.bucket, input.key, params, headers)

	// 发送POST请求
	resp, err = o.httpClient.Do(req)
	// 解析Response到output中
	err = parseResponseToBaseModel(resp, output, true, true)
	if nil != err {
		output = nil
		return
	}
	parseCompleteMultipartUploadOutput(output)
	if output.encodingType == "url" {
		err = decodeCompleteMultipartUploadOutput(output)
		if err != nil {
			output = nil
		}
	}

	return
}

func (o obsClient) abortMultipartUpload(input *abortMultipartUploadInput, extensions ...extensionOptions) (output *baseModel, err error) {
	var (
		redirectURL string
		requestURL  string
		req         *http.Request
		resp        *http.Response
		data        interface{}
	)

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
		return nil, errors.New("abortMultipartUploadInput is nil")
	}
	if input.uploadId == "" {
		return nil, errors.New("uploadId is empty")
	}
	if strings.TrimSpace(input.bucket) == "" && !o.conf.cname {
		err = errors.New("bucket is empty")
		return
	}
	if strings.TrimSpace(input.key) == "" {
		err = errors.New("key is empty")
		return
	}

	//准备参数、请求头、数据
	params, headers, data, err = input.trans(true)
	headers = prepareHeaders(headers, false, true)

	_data, _err := prepareData(headers, data)
	if _err != nil {
		return nil, _err
	}

	req, err = o.getRequest(redirectURL, requestURL, false, _data, "DELETE", input.bucket, input.key, params, headers)

	// 发送POST请求
	resp, err = o.httpClient.Do(req)
	// 解析Response到output中
	err = parseResponseToBaseModel(resp, output, true, true)

	return
}

func (o obsClient) deleteObject(input *deleteObjectInput, extensions ...extensionOptions) (output *deleteObjectOutput, err error) {
	var (
		redirectURL string
		requestURL  string
		req         *http.Request
		resp        *http.Response
		data        interface{}
	)

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
		return nil, errors.New("abortMultipartUploadInput is nil")
	}
	if strings.TrimSpace(input.bucket) == "" && !o.conf.cname {
		err = errors.New("bucket is empty")
		return
	}
	if strings.TrimSpace(input.key) == "" {
		err = errors.New("key is empty")
		return
	}

	//准备参数、请求头、数据
	params, headers, data, err = input.trans(true)
	headers = prepareHeaders(headers, false, true)

	_data, _err := prepareData(headers, data)
	if _err != nil {
		return nil, _err
	}

	req, err = o.getRequest(redirectURL, requestURL, false, _data, "DELETE", input.bucket, input.key, params, headers)

	// 发送POST请求
	resp, err = o.httpClient.Do(req)
	// 解析Response到output中
	err = parseResponseToBaseModel(resp, output, true, true)
	if nil != err {
		output = nil
		return
	}
	parseDeleteObjectOutput(output)

	return
}
