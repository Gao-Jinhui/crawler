package spider

import (
	"crawler/pkg/proxy"
	"go.uber.org/zap"
	"time"
)

type BrowserFetcherOptions struct {
	Logger  *zap.Logger
	Proxy   proxy.ProxyFunc
	Timeout time.Duration
}

type BrowserFetcherOption func(ops *BrowserFetcherOptions)

var DefaultBrowserOptions = BrowserFetcherOptions{
	Proxy: nil,
}

func WithProxy(proxy proxy.ProxyFunc) BrowserFetcherOption {
	return func(ops *BrowserFetcherOptions) {
		ops.Proxy = proxy
	}
}

func WithTimeout(time time.Duration) BrowserFetcherOption {
	return func(ops *BrowserFetcherOptions) {
		ops.Timeout = time
	}
}
