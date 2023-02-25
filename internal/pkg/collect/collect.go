package collect

import (
	"bufio"
	"crawler/pkg/format"
	"crawler/pkg/useragent"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
)

type Fetcher interface {
	Get(resq *Request) ([]byte, error)
}

type BrowserFetch struct {
	options
}

func NewBrowserFetch(opts ...Option) *BrowserFetch {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	bf := &BrowserFetch{}
	bf.options = options
	return bf
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
	if len(request.Task.Cookie) > 0 {
		req.Header.Set("Cookie", request.Task.Cookie)
	}
	req.Header.Set("User-Agent", useragent.GenerateRandomUA())
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.50")

	resp, err := client.Do(req)
	
	if err != nil {
		b.Logger.Error("fetch failed",
			zap.Error(err),
		)
		return nil, err
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := format.DeterminEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}
