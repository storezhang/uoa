package uoa

import (
	`bytes`
	`errors`
	`io`
	`net`
	`net/http`
	`net/url`
	`strings`
	`time`
)

// 准备请求头
func prepareHeaders(headers map[string][]string, meta bool, isObs bool) map[string][]string {
	_headers := make(map[string][]string, len(headers))
	if headers != nil {
		for key, value := range headers {
			key = strings.TrimSpace(key)
			if key == "" {
				continue
			}
			_key := strings.ToLower(key)
			if _, ok := allowedRequestHTTPHeaderMetadataNames[_key]; !ok && !strings.HasPrefix(key, HEADER_PREFIX) && !strings.HasPrefix(key, HEADER_PREFIX_OBS) {
				if !meta {
					continue
				}
				if !isObs {
					_key = HEADER_PREFIX_META + _key
				} else {
					_key = HEADER_PREFIX_META_OBS + _key
				}
			} else {
				_key = key
			}
			_headers[_key] = value
		}
	}

	return _headers
}

// 准备请求数据, 输入参数中的data为任意类型
func prepareData(headers map[string][]string, data interface{}) (io.Reader, error) {
	var _data io.Reader
	if data != nil {
		if dataStr, ok := data.(string); ok {
			headers["Content-Length"] = []string{IntToString(len(dataStr))}
			_data = strings.NewReader(dataStr)
		} else if dataByte, ok := data.([]byte); ok {
			headers["Content-Length"] = []string{IntToString(len(dataByte))}
			_data = bytes.NewReader(dataByte)
		} else if dataReader, ok := data.(io.Reader); ok {
			_data = dataReader
		} else {
			return nil, errors.New("Data is not a valid io.Reader")
		}
	}

	return _data, nil
}

// 获取HTTP请求的 Request
func (o ObsClient) getRequest(redirectURL, requestURL string, redirectFlag bool, _data io.Reader, method,
	bucketName, objectKey string, params map[string]string, headers map[string][]string) (*http.Request, error) {
	if redirectURL != "" {
		if !redirectFlag {
			parsedRedirectURL, err := url.Parse(redirectURL)
			if err != nil {
				return nil, err
			}
			requestURL, err = o.doAuth(method, bucketName, objectKey, params, headers, parsedRedirectURL.Host)
			if err != nil {
				return nil, err
			}
			if parsedRequestURL, err := url.Parse(requestURL); err != nil {
				return nil, err
			} else if parsedRequestURL.RawQuery != "" && parsedRedirectURL.RawQuery == "" {
				redirectURL += "?" + parsedRequestURL.RawQuery
			}
		}
		requestURL = redirectURL
	} else {
		var err error
		requestURL, err = o.doAuth(method, bucketName, objectKey, params, headers, "")
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, requestURL, _data)
	if o.conf.ctx != nil {
		req = req.WithContext(o.conf.ctx)
	}
	if err != nil {
		return nil, err
	}

	return req, nil
}

type connDelegate struct {
	conn          net.Conn
	socketTimeout time.Duration
	finalTimeout  time.Duration
}

func getConnDelegate(conn net.Conn, socketTimeout int, finalTimeout int) *connDelegate {
	return &connDelegate{
		conn:          conn,
		socketTimeout: time.Second * time.Duration(socketTimeout),
		finalTimeout:  time.Second * time.Duration(finalTimeout),
	}
}

func (delegate *connDelegate) Read(b []byte) (n int, err error) {
	err = delegate.SetReadDeadline(time.Now().Add(delegate.socketTimeout))
	if nil != err {
		return
	}
	n, err = delegate.conn.Read(b)
	err = delegate.SetReadDeadline(time.Now().Add(delegate.finalTimeout))
	if nil != err {
		return
	}

	return
}

func (delegate *connDelegate) Write(b []byte) (n int, err error) {
	err = delegate.SetWriteDeadline(time.Now().Add(delegate.socketTimeout))
	if nil != err {
		return
	}

	n, err = delegate.conn.Write(b)
	err = delegate.SetWriteDeadline(time.Now().Add(delegate.finalTimeout))
	if nil != err {
		return
	}

	err = delegate.SetReadDeadline(time.Now().Add(delegate.finalTimeout))

	return
}

func (delegate *connDelegate) Close() error {
	return delegate.conn.Close()
}

func (delegate *connDelegate) LocalAddr() net.Addr {
	return delegate.conn.LocalAddr()
}

func (delegate *connDelegate) RemoteAddr() net.Addr {
	return delegate.conn.RemoteAddr()
}

func (delegate *connDelegate) SetDeadline(t time.Time) error {
	return delegate.conn.SetDeadline(t)
}

func (delegate *connDelegate) SetReadDeadline(t time.Time) error {
	return delegate.conn.SetReadDeadline(t)
}

func (delegate *connDelegate) SetWriteDeadline(t time.Time) error {
	return delegate.conn.SetWriteDeadline(t)
}
