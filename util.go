package uoa

import (
	`crypto/hmac`
	`crypto/md5`
	`crypto/sha1`
	`crypto/sha256`
	`encoding/base64`
	`encoding/hex`
	`encoding/json`
	`net/url`
	`regexp`
	`strconv`
	`strings`
	`time`
)

var regex = regexp.MustCompile("^[\u4e00-\u9fa5]$")
var ipRegex = regexp.MustCompile("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$")
var v4AuthRegex = regexp.MustCompile("Credential=(.+?),SignedHeaders=(.+?),Signature=.+")
var regionRegex = regexp.MustCompile(".+/\\d+/(.+?)/.+")

func StringContains(src string, subStr string, subTranscoding string) string {
	return strings.Replace(src, subStr, subTranscoding, -1)
}

func StringToInt(value string, def int) int {
	ret, err := strconv.Atoi(value)
	if err != nil {
		ret = def
	}

	return ret
}

func StringToInt64(value string, def int64) int64 {
	ret, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		ret = def
	}

	return ret
}

func IntToString(value int) string {
	return strconv.Itoa(value)
}

func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano() / 1000000
}

func FormatUtcNow(format string) string {
	return time.Now().UTC().Format(format)
}

func Md5(value []byte) ([]byte, error) {
	m := md5.New()
	_, err := m.Write(value)
	if err != nil {
		return nil, err
	}

	return m.Sum(nil), nil
}

func HmacSha1(key, value []byte) []byte {
	mac := hmac.New(sha1.New, key)
	_, err := mac.Write(value)
	if err != nil {
		return nil
	}

	return mac.Sum(nil)
}

func HmacSha256(key, value []byte) []byte {
	mac := hmac.New(sha256.New, key)
	_, err := mac.Write(value)
	if err != nil {
		return nil
	}

	return mac.Sum(nil)
}

func Base64Encode(value []byte) string {
	return base64.StdEncoding.EncodeToString(value)
}

func Base64Decode(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}

func HexMd5(value []byte) string {
	bytes, _ := Md5(value)

	return Hex(bytes)
}

func Base64Md5(value []byte) string {
	bytes, _ := Md5(value)

	return Base64Encode(bytes)
}

func parseJSON(value []byte, result interface{}) error {
	if len(value) == 0 {
		return nil
	}

	return json.Unmarshal(value, result)
}

func Hex(value []byte) string {
	return hex.EncodeToString(value)
}

func HexSha256(value []byte) string {
	return Hex(Sha256Hash(value))
}

func Sha256Hash(value []byte) []byte {
	hash := sha256.New()
	_, err := hash.Write(value)
	if err != nil {
		return nil
	}

	return hash.Sum(nil)
}

func IsIP(value string) bool {
	return ipRegex.MatchString(value)
}

func isPathStyle(headers map[string][]string, bucketName string) bool {
	if receivedHost, ok := headers[HEADER_HOST]; ok && len(receivedHost) > 0 && !strings.HasPrefix(receivedHost[0], bucketName+".") {
		return true
	}

	return false
}

// 将时间字符串格式转换成 RFC1123 格式
func FormatUtcToRfc1123(t time.Time) string {
	ret := t.UTC().Format(time.RFC1123)
	return ret[:strings.LastIndex(ret, "UTC")] + "GMT"
}

// URL转码，将中文字符转换为国际码
func UrlEncode(value string, chineseOnly bool) string {
	if chineseOnly {
		values := make([]string, 0, len(value))
		for _, val := range value {
			_value := string(val)
			if regex.MatchString(_value) {
				_value = url.QueryEscape(_value)
			}
			values = append(values, _value)
		}
		return strings.Join(values, "")
	}

	return url.QueryEscape(value)
}
