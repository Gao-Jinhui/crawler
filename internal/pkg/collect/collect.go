package collect

import (
	"bufio"
	"crawler/pkg/format"
	"crawler/pkg/proxy"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"time"
)

type Fetcher interface {
	Get(resq *Request) ([]byte, error)
}

type BrowserFetch struct {
	Timeout time.Duration
	Proxy   proxy.ProxyFunc
	Logger  *zap.Logger
}

func NewBrowserFetch(logger *zap.Logger, proxy proxy.ProxyFunc) *BrowserFetch {
	return &BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   proxy,
		Logger:  logger,
	}
}

// Get 模拟浏览器访问
func (b *BrowserFetch) Get(request *Request) ([]byte, error) {

	client := &http.Client{
		Timeout: b.Timeout,
	}
	if b.Proxy != nil {
		transport := http.DefaultTransport.(*http.Transport)
		transport.Proxy = b.Proxy
		client.Transport = transport
	}
	req, err := http.NewRequest("GET", request.Url, nil)
	if err != nil {
		return nil, fmt.Errorf("get url failed:%v", err)
	}
	if len(request.Cookie) > 0 {
		req.Header.Set("Cookie", request.Cookie)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := format.DeterminEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}
