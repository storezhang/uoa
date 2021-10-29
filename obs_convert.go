package uoa

import (
	`encoding/xml`
	`fmt`
	`io/ioutil`
	`net/http`
	`net/url`
	`reflect`
	`strings`
)

func parseSseHeader(responseHeaders map[string][]string) (sseHeader isseHeader) {
	if ret, ok := responseHeaders[HEADER_SSEC_ENCRYPTION]; ok {
		sseCHeader := sseCHeader{encryption: ret[0]}
		if ret, ok = responseHeaders[HEADER_SSEC_KEY_MD5]; ok {
			sseCHeader.keyMD5 = ret[0]
		}
		sseHeader = sseCHeader
	} else if ret, ok := responseHeaders[HEADER_SSEKMS_ENCRYPTION]; ok {
		sseKmsHeader := sseKmsHeader{encryption: ret[0]}
		if ret, ok = responseHeaders[HEADER_SSEKMS_KEY]; ok {
			sseKmsHeader.key = ret[0]
		} else if ret, ok = responseHeaders[HEADER_SSEKMS_ENCRYPT_KEY_OBS]; ok {
			sseKmsHeader.key = ret[0]
		}
		sseHeader = sseKmsHeader
	}

	return
}

func parseInitiateMultipartUploadOutput(output *initiateMultipartUploadOutput) {
	output.sseHeader = parseSseHeader(output.responseHeaders)
}

func decodeInitiateMultipartUploadOutput(output *initiateMultipartUploadOutput) (err error) {
	output.key, err = url.QueryUnescape(output.key)
	return
}

func parseResponseToBaseModel(resp *http.Response, baseModel ibaseModel, xmlResult bool, isObs bool) (err error) {
	readCloser, ok := baseModel.(ireadCloser)
	if !ok {
		defer func() {
			_ = resp.Body.Close()
		}()
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil && len(body) > 0 {
			if xmlResult {
				err = parseXml(body, baseModel)
			} else {
				s := reflect.TypeOf(baseModel).Elem()
				if reflect.TypeOf(baseModel).Elem().Name() == "GetBucketPolicyOutput" {
					parseBucketPolicyOutput(s, baseModel, body)
				} else {
					err = parseJSON(body, baseModel)
				}
			}
			if err != nil {
				return
			}
		}
	} else {
		readCloser.setReadCloser(resp.Body)
	}

	baseModel.setStatusCode(resp.StatusCode)
	responseHeaders := cleanHeaderPrefix(resp.Header)
	baseModel.setResponseHeaders(responseHeaders)
	if values, ok := responseHeaders[HEADER_REQUEST_ID]; ok {
		baseModel.setRequestID(values[0])
	}

	return
}

func parseXml(value []byte, result interface{}) error {
	if len(value) == 0 {
		return nil
	}

	return xml.Unmarshal(value, result)
}

func parseBucketPolicyOutput(s reflect.Type, baseModel ibaseModel, body []byte) {
	for i := 0; i < s.NumField(); i++ {
		if s.Field(i).Tag == "json:\"body\"" {
			reflect.ValueOf(baseModel).Elem().FieldByName(s.Field(i).Name).SetString(string(body))
			break
		}
	}
}

func cleanHeaderPrefix(header http.Header) map[string][]string {
	responseHeaders := make(map[string][]string)
	for key, value := range header {
		if len(value) > 0 {
			key = strings.ToLower(key)
			if strings.HasPrefix(key, HEADER_PREFIX) || strings.HasPrefix(key, HEADER_PREFIX_OBS) {
				key = key[len(HEADER_PREFIX):]
			}
			responseHeaders[key] = value
		}
	}

	return responseHeaders
}

func convertCompleteMultipartUploadInputToXml(input completeMultipartUploadInput, returnMd5 bool) (data string, md5 string) {
	xml := make([]string, 0, 2+len(input.parts)*4)
	xml = append(xml, "<CompleteMultipartUpload>")
	for _, part := range input.parts {
		xml = append(xml, "<Part>")
		xml = append(xml, fmt.Sprintf("<PartNumber>%d</PartNumber>", part.partNumber))
		xml = append(xml, fmt.Sprintf("<ETag>%s</ETag>", part.eTag))
		xml = append(xml, "</Part>")
	}
	xml = append(xml, "</CompleteMultipartUpload>")
	data = strings.Join(xml, "")
	if returnMd5 {
		md5 = base64Md5([]byte(data))
	}

	return
}

func parseCompleteMultipartUploadOutput(output *completeMultipartUploadOutput) {
	output.sseHeader = parseSseHeader(output.responseHeaders)
	if ret, ok := output.responseHeaders[HEADER_VERSION_ID]; ok {
		output.versionId = ret[0]
	}
}

func decodeCompleteMultipartUploadOutput(output *completeMultipartUploadOutput) (err error) {
	output.key, err = url.QueryUnescape(output.key)
	return
}

func parseDeleteObjectOutput(output *deleteObjectOutput) {
	if versionID, ok := output.responseHeaders[HEADER_VERSION_ID]; ok {
		output.versionId = versionID[0]
	}

	if deleteMarker, ok := output.responseHeaders[HEADER_DELETE_MARKER]; ok {
		output.deleteMarker = deleteMarker[0] == "true"
	}
}
