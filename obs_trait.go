package uoa

import "strings"

func (input ObjectOperationInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string)
	params = make(map[string]string)
	if acl := string(input.ACL); acl != "" {
		setHeaders(headers, HEADER_ACL, []string{acl}, isObs)
	}
	input.prepareGrantHeaders(headers)
	if storageClass := string(input.StorageClass); storageClass != "" {
		if !isObs {
			if storageClass == string(StorageClassWarm) {
				storageClass = string(storageClassStandardIA)
			} else if storageClass == string(StorageClassCold) {
				storageClass = string(storageClassGlacier)
			}
		}
		setHeaders(headers, HEADER_STORAGE_CLASS2, []string{storageClass}, isObs)
	}
	if input.WebsiteRedirectLocation != "" {
		setHeaders(headers, HEADER_WEBSITE_REDIRECT_LOCATION, []string{input.WebsiteRedirectLocation}, isObs)

	}
	setSseHeader(headers, input.SseHeader, false, isObs)
	if input.Expires != 0 {
		setHeaders(headers, HEADER_EXPIRES, []string{Int64ToString(input.Expires)}, true)
	}
	if input.Metadata != nil {
		for key, value := range input.Metadata {
			key = strings.TrimSpace(key)
			setHeadersNext(headers, HEADER_PREFIX_META_OBS+key, HEADER_PREFIX_META+key, []string{value}, isObs)
		}
	}

	return
}

func (input HeadObjectInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = make(map[string]string)
	if input.VersionId != "" {
		params[PARAM_VERSION_ID] = input.VersionId
	}
	return
}

func (input InitiateMultipartUploadInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.ObjectOperationInput.trans(isObs)
	if err != nil {
		return
	}
	if input.ContentType != "" {
		headers[HEADER_CONTENT_TYPE_CAML] = []string{input.ContentType}
	}
	params[string("uploads")] = ""
	if input.EncodingType != "" {
		params["encoding-type"] = input.EncodingType
	}

	return
}

func (input CompleteMultipartUploadInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{"uploadId": input.UploadId}
	if input.EncodingType != "" {
		params["encoding-type"] = input.EncodingType
	}
	data, _ = ConvertCompleteMultipartUploadInputToXml(input, false)

	return
}

func (input AbortMultipartUploadInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{"uploadId": input.UploadId}
	return
}

func (input DeleteObjectInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = make(map[string]string)
	if input.VersionId != "" {
		params[PARAM_VERSION_ID] = input.VersionId
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

func (input ObjectOperationInput) prepareGrantHeaders(headers map[string][]string) {
	if GrantReadID := input.GrantReadId; GrantReadID != "" {
		setHeaders(headers, HEADER_GRANT_READ_OBS, []string{GrantReadID}, true)
	}
	if GrantReadAcpID := input.GrantReadAcpId; GrantReadAcpID != "" {
		setHeaders(headers, HEADER_GRANT_READ_ACP_OBS, []string{GrantReadAcpID}, true)
	}
	if GrantWriteAcpID := input.GrantWriteAcpId; GrantWriteAcpID != "" {
		setHeaders(headers, HEADER_GRANT_WRITE_ACP_OBS, []string{GrantWriteAcpID}, true)
	}
	if GrantFullControlID := input.GrantFullControlId; GrantFullControlID != "" {
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

func setSseHeader(headers map[string][]string, sseHeader ISseHeader, sseCOnly bool, isObs bool) {
	if sseHeader != nil {
		if sseCHeader, ok := sseHeader.(SseCHeader); ok {
			setHeaders(headers, HEADER_SSEC_ENCRYPTION, []string{sseCHeader.GetEncryption()}, isObs)
			setHeaders(headers, HEADER_SSEC_KEY, []string{sseCHeader.GetKey()}, isObs)
			setHeaders(headers, HEADER_SSEC_KEY_MD5, []string{sseCHeader.GetKeyMD5()}, isObs)
		} else if sseKmsHeader, ok := sseHeader.(SseKmsHeader); !sseCOnly && ok {
			sseKmsHeader.isObs = isObs
			setHeaders(headers, HEADER_SSEKMS_ENCRYPTION, []string{sseKmsHeader.GetEncryption()}, isObs)
			if sseKmsHeader.GetKey() != "" {
				setHeadersNext(headers, HEADER_SSEKMS_KEY_OBS, HEADER_SSEKMS_KEY_AMZ, []string{sseKmsHeader.GetKey()}, isObs)
			}
		}
	}
}

func (header SseCHeader) GetKeyMD5() string {
	if header.KeyMD5 != "" {
		return header.KeyMD5
	}
	if ret, err := Base64Decode(header.GetKey()); err == nil {
		return Base64Md5(ret)
	}

	return ""
}
