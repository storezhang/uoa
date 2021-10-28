package uoa

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

func parseSseHeader(responseHeaders map[string][]string) (sseHeader ISseHeader) {
	if ret, ok := responseHeaders[HEADER_SSEC_ENCRYPTION]; ok {
		sseCHeader := SseCHeader{Encryption: ret[0]}
		if ret, ok = responseHeaders[HEADER_SSEC_KEY_MD5]; ok {
			sseCHeader.KeyMD5 = ret[0]
		}
		sseHeader = sseCHeader
	} else if ret, ok := responseHeaders[HEADER_SSEKMS_ENCRYPTION]; ok {
		sseKmsHeader := SseKmsHeader{Encryption: ret[0]}
		if ret, ok = responseHeaders[HEADER_SSEKMS_KEY]; ok {
			sseKmsHeader.Key = ret[0]
		} else if ret, ok = responseHeaders[HEADER_SSEKMS_ENCRYPT_KEY_OBS]; ok {
			sseKmsHeader.Key = ret[0]
		}
		sseHeader = sseKmsHeader
	}

	return
}

func ParseInitiateMultipartUploadOutput(output *InitiateMultipartUploadOutput) {
	output.SseHeader = parseSseHeader(output.ResponseHeaders)
}

func decodeInitiateMultipartUploadOutput(output *InitiateMultipartUploadOutput) (err error) {
	output.Key, err = url.QueryUnescape(output.Key)
	return
}

func ParseResponseToBaseModel(resp *http.Response, baseModel IBaseModel, xmlResult bool, isObs bool) (err error) {
	readCloser, ok := baseModel.(IReadCloser)
	if !ok {
		defer func() {
			_ = resp.Body.Close()
		}()
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil && len(body) > 0 {
			if xmlResult {
				err = ParseXml(body, baseModel)
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

func ParseXml(value []byte, result interface{}) error {
	if len(value) == 0 {
		return nil
	}

	return xml.Unmarshal(value, result)
}

func parseBucketPolicyOutput(s reflect.Type, baseModel IBaseModel, body []byte) {
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

func ConvertCompleteMultipartUploadInputToXml(input CompleteMultipartUploadInput, returnMd5 bool) (data string, md5 string) {
	xml := make([]string, 0, 2+len(input.Parts)*4)
	xml = append(xml, "<CompleteMultipartUpload>")
	for _, part := range input.Parts {
		xml = append(xml, "<Part>")
		xml = append(xml, fmt.Sprintf("<PartNumber>%d</PartNumber>", part.PartNumber))
		xml = append(xml, fmt.Sprintf("<ETag>%s</ETag>", part.ETag))
		xml = append(xml, "</Part>")
	}
	xml = append(xml, "</CompleteMultipartUpload>")
	data = strings.Join(xml, "")
	if returnMd5 {
		md5 = Base64Md5([]byte(data))
	}

	return
}

func ParseCompleteMultipartUploadOutput(output *CompleteMultipartUploadOutput) {
	output.SseHeader = parseSseHeader(output.ResponseHeaders)
	if ret, ok := output.ResponseHeaders[HEADER_VERSION_ID]; ok {
		output.VersionId = ret[0]
	}
}

func decodeCompleteMultipartUploadOutput(output *CompleteMultipartUploadOutput) (err error) {
	output.Key, err = url.QueryUnescape(output.Key)
	return
}

func ParseDeleteObjectOutput(output *DeleteObjectOutput) {
	if versionID, ok := output.ResponseHeaders[HEADER_VERSION_ID]; ok {
		output.VersionId = versionID[0]
	}

	if deleteMarker, ok := output.ResponseHeaders[HEADER_DELETE_MARKER]; ok {
		output.DeleteMarker = deleteMarker[0] == "true"
	}
}
