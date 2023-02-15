package collect

import (
	"crawler/pkg/proxy"
	"go.uber.org/zap"
	"time"
)

type options struct {
	Logger  *zap.Logger
	Proxy   proxy.ProxyFunc
	Timeout time.Duration
}

type Option func(ops *options)

var defaultOptions = options{
	Logger: zap.NewNop(),
	Proxy:  nil,
}

func WithLogger(logger *zap.Logger) Option {
	return func(ops *options) {
		ops.Logger = logger
	}
}

func WithProxy(proxy proxy.ProxyFunc) Option {
	return func(ops *options) {
		ops.Proxy = proxy
	}
}

func WithTimeout(time time.Duration) Option {
	return func(ops *options) {
		ops.Timeout = time
	}
}
