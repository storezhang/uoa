package uoa

import `strings`

func (input objectOperationInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string)
	params = make(map[string]string)
	if acl := string(input.acl); acl != "" {
		setHeaders(headers, HEADER_ACL, []string{acl}, isObs)
	}
	input.prepareGrantHeaders(headers)
	if storageClass := string(input.storageClass); storageClass != "" {
		if !isObs {
			if storageClass == string(StorageClassWarm) {
				storageClass = string(storageClassStandardIA)
			} else if storageClass == string(StorageClassCold) {
				storageClass = string(storageClassGlacier)
			}
		}
		setHeaders(headers, HEADER_STORAGE_CLASS2, []string{storageClass}, isObs)
	}
	if input.websiteRedirectLocation != "" {
		setHeaders(headers, HEADER_WEBSITE_REDIRECT_LOCATION, []string{input.websiteRedirectLocation}, isObs)

	}
	setSseHeader(headers, input.sseHeader, false, isObs)
	if input.expires != 0 {
		setHeaders(headers, HEADER_EXPIRES, []string{int64ToString(input.expires)}, true)
	}
	if input.metadata != nil {
		for key, value := range input.metadata {
			key = strings.TrimSpace(key)
			setHeadersNext(headers, HEADER_PREFIX_META_OBS+key, HEADER_PREFIX_META+key, []string{value}, isObs)
		}
	}

	return
}

func (input headObjectInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = make(map[string]string)
	if input.versionId != "" {
		params[PARAM_VERSION_ID] = input.versionId
	}
	return
}

func (input initiateMultipartUploadInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.objectOperationInput.trans(isObs)
	if err != nil {
		return
	}
	if input.contentType != "" {
		headers[HEADER_CONTENT_TYPE_CAML] = []string{input.contentType}
	}
	params[string("uploads")] = ""
	if input.encodingType != "" {
		params["encoding-type"] = input.encodingType
	}

	return
}

func (input completeMultipartUploadInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{"uploadId": input.uploadId}
	if input.encodingType != "" {
		params["encoding-type"] = input.encodingType
	}
	data, _ = convertCompleteMultipartUploadInputToXml(input, false)

	return
}

func (input abortMultipartUploadInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{"uploadId": input.uploadId}
	return
}

func (input deleteObjectInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = make(map[string]string)
	if input.versionId != "" {
		params[PARAM_VERSION_ID] = input.versionId
	}
	return
}

func setHeaders(headers map[string][]string, header string, headerValue []string, isObs bool) {
	if isObs {
		header = HEADER_PREFIX_OBS + header
		headers[header] = headerValue
	} else {
		header = HEADER_PREFIX + header
		headers[header] = headerValue
	}
}

func (input objectOperationInput) prepareGrantHeaders(headers map[string][]string) {
	if GrantReadID := input.grantReadId; GrantReadID != "" {
		setHeaders(headers, HEADER_GRANT_READ_OBS, []string{GrantReadID}, true)
	}
	if GrantReadAcpID := input.grantReadAcpId; GrantReadAcpID != "" {
		setHeaders(headers, HEADER_GRANT_READ_ACP_OBS, []string{GrantReadAcpID}, true)
	}
	if GrantWriteAcpID := input.grantWriteAcpId; GrantWriteAcpID != "" {
		setHeaders(headers, HEADER_GRANT_WRITE_ACP_OBS, []string{GrantWriteAcpID}, true)
	}
	if GrantFullControlID := input.grantFullControlId; GrantFullControlID != "" {
		setHeaders(headers, HEADER_GRANT_FULL_CONTROL_OBS, []string{GrantFullControlID}, true)
	}
}

func setHeadersNext(headers map[string][]string, header string, headerNext string, headerValue []string, isObs bool) {
	if isObs {
		headers[header] = headerValue
	} else {
		headers[headerNext] = headerValue
	}
}

func setSseHeader(headers map[string][]string, sseHeader isseHeader, sseCOnly bool, isObs bool) {
	if sseHeader != nil {
		if sseCHeader, ok := sseHeader.(sseCHeader); ok {
			setHeaders(headers, HEADER_SSEC_ENCRYPTION, []string{sseCHeader.getEncryption()}, isObs)
			setHeaders(headers, HEADER_SSEC_KEY, []string{sseCHeader.getKey()}, isObs)
			setHeaders(headers, HEADER_SSEC_KEY_MD5, []string{sseCHeader.GetKeyMD5()}, isObs)
		} else if sseKmsHeader, ok := sseHeader.(sseKmsHeader); !sseCOnly && ok {
			sseKmsHeader.isObs = isObs
			setHeaders(headers, HEADER_SSEKMS_ENCRYPTION, []string{sseKmsHeader.getEncryption()}, isObs)
			if sseKmsHeader.getKey() != "" {
				setHeadersNext(headers, HEADER_SSEKMS_KEY_OBS, HEADER_SSEKMS_KEY_AMZ, []string{sseKmsHeader.getKey()}, isObs)
			}
		}
	}
}

func (header sseCHeader) GetKeyMD5() string {
	if header.keyMD5 != "" {
		return header.keyMD5
	}
	if ret, err := base64Decode(header.getKey()); err == nil {
		return base64Md5(ret)
	}

	return ""
}
