package uoa

import (
	`context`
	`crypto/tls`
	`crypto/x509`
	`errors`
	`fmt`
	`net`
	`net/http`
	`net/url`
	`sort`
	`strconv`
	`strings`
	`time`
)

type urlHolder struct {
	scheme string
	host   string
	port   int
}

// obs 配置结构体
type config struct {
	securityProviders []securityProvider
	urlHolder         *urlHolder
	pathStyle         bool
	cname             bool
	sslVerify         bool
	endPoint          string
	signature         SignatureType
	region            string
	connectTimeout    int
	socketTimeout     int
	headerTimeout     int
	idleConnTimeout   int
	finalTimeout      int
	maxRetryCount     int
	proxyURL          string
	maxConnsPerHost   int
	pemCerts          []byte
	transport         *http.Transport
	ctx               context.Context
	maxRedirectCount  int
	userAgent         string
	enableCompression bool
}

func (conf config) String() string {
	return fmt.Sprintf("[endPoint:%s, signature:%s, pathStyle:%v, region:%s"+
		"\nconnectTimeout:%d, socketTimeout:%dheaderTimeout:%d, idleConnTimeout:%d"+
		"\nmaxRetryCount:%d, maxConnsPerHost:%d, sslVerify:%v, maxRedirectCount:%d]",
		conf.endPoint, conf.signature, conf.pathStyle, conf.region,
		conf.connectTimeout, conf.socketTimeout, conf.headerTimeout, conf.idleConnTimeout,
		conf.maxRetryCount, conf.maxConnsPerHost, conf.sslVerify, conf.maxRedirectCount)
}

type configuror func(conf *config)

func WithSecurityProviders(sps ...securityProvider) configuror {
	return func(conf *config) {
		for _, sp := range sps {
			if nil != sp {
				conf.securityProviders = append(conf.securityProviders, sp)
			}
		}
	}
}

func WithSslVerify(sslVerify bool) configuror {
	return WithSslVerifyAndPemCerts(sslVerify, nil)
}

func WithSslVerifyAndPemCerts(sslVerify bool, pemCerts []byte) configuror {
	return func(conf *config) {
		conf.sslVerify = sslVerify
		conf.pemCerts = pemCerts
	}
}

func WithHeaderTimeOut(headerTimeout int) configuror {
	return func(conf *config) {
		conf.headerTimeout = headerTimeout
	}
}

func WithProxyURL(proxyURL string) configuror {
	return func(conf *config) {
		conf.proxyURL = proxyURL
	}
}

func WithMaxConnects(maxConnPerHost int) configuror {
	return func(conf *config) {
		conf.maxConnsPerHost = maxConnPerHost
	}
}

func WithPathStyle(pathStyle bool) configuror {
	return func(conf *config) {
		conf.pathStyle = pathStyle
	}
}

func WithSignature(signature SignatureType) configuror {
	return func(conf *config) {
		conf.signature = signature
	}
}

func WithRegion(region string) configuror {
	return func(conf *config) {
		conf.region = region
	}
}

func WithConnectTimeout(connectTimeout int) configuror {
	return func(conf *config) {
		conf.connectTimeout = connectTimeout
	}
}

func WithSocketTimeout(socketTimeout int) configuror {
	return func(conf *config) {
		conf.socketTimeout = socketTimeout
	}
}

func WithIdleConnTimeout(idleConnTimeout int) configuror {
	return func(conf *config) {
		conf.idleConnTimeout = idleConnTimeout
	}
}

func WithMaxRetryCount(maxRetryCount int) configuror {
	return func(conf *config) {
		conf.maxRetryCount = maxRetryCount
	}
}

func WithSecurityToken(securityToken string) configuror {
	return func(conf *config) {
		for _, sp := range conf.securityProviders {
			if bsp, ok := sp.(*BasicSecurityProvider); ok {
				sh := bsp.getSecurity()
				bsp.refresh(sh.accessKey, sh.accessKey, securityToken)
				break
			}
		}
	}
}

func WithHttpTransport(transport *http.Transport) configuror {
	return func(conf *config) {
		conf.transport = transport
	}
}

func WithRequestContext(ctx context.Context) configuror {
	return func(conf *config) {
		conf.ctx = ctx
	}
}

func WithCustomDomainName(cname bool) configuror {
	return func(conf *config) {
		conf.cname = cname
	}
}

func WithMaxRedirectCount(maxRedirectCount int) configuror {
	return func(conf *config) {
		conf.maxRedirectCount = maxRedirectCount
	}
}

func WithUserAgent(userAgent string) configuror {
	return func(conf *config) {
		conf.userAgent = userAgent
	}
}

func WithEnableCompression(enableCompression bool) configuror {
	return func(conf *config) {
		conf.enableCompression = enableCompression
	}
}

func (conf *config) prepareConfig() {
	if conf.connectTimeout <= 0 {
		conf.connectTimeout = DEFAULT_CONNECT_TIMEOUT
	}

	if conf.socketTimeout <= 0 {
		conf.socketTimeout = DEFAULT_SOCKET_TIMEOUT
	}

	conf.finalTimeout = conf.socketTimeout * 10

	if conf.headerTimeout <= 0 {
		conf.headerTimeout = DEFAULT_HEADER_TIMEOUT
	}

	if conf.idleConnTimeout < 0 {
		conf.idleConnTimeout = DEFAULT_IDLE_CONN_TIMEOUT
	}

	if conf.maxRetryCount < 0 {
		conf.maxRetryCount = DEFAULT_MAX_RETRY_COUNT
	}

	if conf.maxConnsPerHost <= 0 {
		conf.maxConnsPerHost = DEFAULT_MAX_CONN_PER_HOST
	}

	if conf.maxRedirectCount < 0 {
		conf.maxRedirectCount = DEFAULT_MAX_REDIRECT_COUNT
	}

	if conf.pathStyle && conf.signature == SignatureObs {
		conf.signature = SignatureV2
	}
}

func (conf *config) initConfigWithDefault() error {
	conf.endPoint = strings.TrimSpace(conf.endPoint)
	if conf.endPoint == "" {
		return errors.New("endpoint is not set")
	}

	if index := strings.Index(conf.endPoint, "?"); index > 0 {
		conf.endPoint = conf.endPoint[:index]
	}

	if strings.LastIndex(conf.endPoint, "/") == len(conf.endPoint)-1 {
		conf.endPoint = conf.endPoint[:len(conf.endPoint)-1]
	}

	if conf.signature == "" {
		conf.signature = DEFAULT_SIGNATURE
	}

	urlHolder := &urlHolder{}
	var address string
	if strings.HasPrefix(conf.endPoint, "https://") {
		urlHolder.scheme = "https"
		address = conf.endPoint[len("https://"):]
	} else if strings.HasPrefix(conf.endPoint, "http") {
		urlHolder.scheme = "http"
		address = conf.endPoint[len("http://"):]
	} else {
		urlHolder.scheme = "https"
		address = conf.endPoint
	}

	addr := strings.Split(address, ":")
	if len(addr) == 2 {
		if port, err := strconv.Atoi(addr[1]); err == nil {
			urlHolder.port = port
		}
	}

	urlHolder.host = addr[0]
	if urlHolder.port == 0 {
		if urlHolder.scheme == "https" {
			urlHolder.port = 443
		} else {
			urlHolder.port = 80
		}
	}

	if IsIP(urlHolder.host) {
		conf.pathStyle = true
	}

	conf.urlHolder = urlHolder
	conf.region = strings.TrimSpace(conf.region)
	if conf.region == "" {
		conf.region = DEFAULT_REGION
	}

	conf.prepareConfig()
	conf.proxyURL = strings.TrimSpace(conf.proxyURL)

	return nil
}

func (conf *config) getTransport() error {
	if conf.transport == nil {
		conf.transport = &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(network, addr, time.Second*time.Duration(conf.connectTimeout))
				if err != nil {
					return nil, err
				}
				return getConnDelegate(conn, conf.socketTimeout, conf.finalTimeout), nil
			},
			MaxIdleConns:          conf.maxConnsPerHost,
			MaxIdleConnsPerHost:   conf.maxConnsPerHost,
			ResponseHeaderTimeout: time.Second * time.Duration(conf.headerTimeout),
			IdleConnTimeout:       time.Second * time.Duration(conf.idleConnTimeout),
		}

		if conf.proxyURL != "" {
			proxyURL, err := url.Parse(conf.proxyURL)
			if err != nil {
				return err
			}
			conf.transport.Proxy = http.ProxyURL(proxyURL)
		}

		tlsConfig := &tls.Config{InsecureSkipVerify: !conf.sslVerify}
		if conf.sslVerify && conf.pemCerts != nil {
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(conf.pemCerts)
			tlsConfig.RootCAs = pool
		}

		conf.transport.TLSClientConfig = tlsConfig
		conf.transport.DisableCompression = !conf.enableCompression
	}

	return nil
}

func checkRedirectFunc(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func DummyQueryEscape(s string) string {
	return s
}

func (conf *config) prepareBaseURL(bucket string) (requestURL string, canonicalizeURL string) {
	urlHolder := conf.urlHolder
	if conf.cname {
		requestURL = fmt.Sprintf("%s://%s:%d", urlHolder.scheme, urlHolder.host, urlHolder.port)
		if conf.signature == "v4" {
			canonicalizeURL = "/"
		} else {
			canonicalizeURL = "/" + urlHolder.host + "/"
		}
	} else {
		if bucket == "" {
			requestURL = fmt.Sprintf("%s://%s:%d", urlHolder.scheme, urlHolder.host, urlHolder.port)
			canonicalizeURL = "/"
		} else {
			if conf.pathStyle {
				requestURL = fmt.Sprintf("%s://%s:%d/%s", urlHolder.scheme, urlHolder.host, urlHolder.port, bucket)
				canonicalizeURL = "/" + bucket
			} else {
				requestURL = fmt.Sprintf("%s://%s.%s:%d", urlHolder.scheme, bucket, urlHolder.host, urlHolder.port)
				if conf.signature == "v2" || conf.signature == "OBS" {
					canonicalizeURL = "/" + bucket + "/"
				} else {
					canonicalizeURL = "/"
				}
			}
		}
	}

	return
}

func (conf *config) prepareObjectKey(escape bool, objectKey string, escapeFunc func(s string) string) (encodeObjectKey string) {
	if escape {
		tempKey := []rune(objectKey)
		result := make([]string, 0, len(tempKey))

		for _, val := range tempKey {
			if string(val) == "/" {
				result = append(result, string(val))
			} else {
				if string(val) == " " {
					result = append(result, url.PathEscape(string(val)))
				} else {
					result = append(result, url.QueryEscape(string(val)))
				}
			}
		}
		encodeObjectKey = strings.Join(result, "")
	} else {
		encodeObjectKey = escapeFunc(objectKey)
	}

	return
}

func (conf *config) prepareEscapeFunc(escape bool) (escapeFunc func(s string) string) {
	if escape {
		return url.QueryEscape
	}

	return DummyQueryEscape
}

func (conf *config) formatUrls(bucket string, objectKey string, params map[string]string, escape bool) (requestURL string, canonicalizeURL string) {
	requestURL, canonicalizeURL = conf.prepareBaseURL(bucket)
	var escapeFunc func(s string) string
	escapeFunc = conf.prepareEscapeFunc(escape)

	if objectKey != "" {
		var encodeObjectKey string
		encodeObjectKey = conf.prepareObjectKey(escape, objectKey, escapeFunc)
		requestURL += "/" + encodeObjectKey
		if !strings.HasPrefix(canonicalizeURL, "/") {
			canonicalizeURL += "/"
		}
		canonicalizeURL += encodeObjectKey
	}

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, strings.TrimSpace(key))
	}
	sort.Strings(keys)
	i := 0
	for index, key := range keys {
		if index == 0 {
			requestURL += "?"
		} else {
			requestURL += "&"
		}
		_key := url.QueryEscape(key)
		requestURL += _key

		_value := params[key]
		if conf.signature == "v4" {
			requestURL += "=" + url.QueryEscape(_value)
		} else {
			if _value != "" {
				requestURL += "=" + url.QueryEscape(_value)
				_value = "=" + _value
			} else {
				_value = ""
			}

			lowerKey := strings.ToLower(key)
			_, ok := allowedResourceParameterNames[lowerKey]
			prefixHeader := HEADER_PREFIX
			isObs := conf.signature == SignatureObs
			if isObs {
				prefixHeader = HEADER_PREFIX_OBS
			}
			ok = ok || strings.HasPrefix(lowerKey, prefixHeader)
			if ok {
				if i == 0 {
					canonicalizeURL += "?"
				} else {
					canonicalizeURL += "&"
				}
				canonicalizeURL += getQueryURL(_key, _value)
				i++
			}
		}
	}

	return
}

func getQueryURL(key, value string) string {
	queryURL := ""
	queryURL += key
	queryURL += value

	return queryURL
}
