package spider

import (
	"crawler/internal/pkg/collector"
	"crawler/internal/pkg/limiter"
	"go.uber.org/zap"
)

type TaskOptions struct {
	Name     string `json:"name"` // 任务名称，应保证唯一性
	Cookie   string `json:"cookie"`
	WaitTime int64  `json:"wait_time"` // 随机休眠时间，秒
	Reload   bool   `json:"reload"`    // 网站是否可以重复爬取
	MaxDepth int64  `json:"max_depth"`
	Fetcher  Fetcher
	Storage  collector.Storage
	Limit    limiter.RateLimiter
	Logger   *zap.Logger
}

var DefaultTaskOptions = TaskOptions{
	Logger:   zap.NewNop(),
	WaitTime: 5,
	Reload:   false,
	MaxDepth: 5,
}

type TaskOption func(opts *TaskOptions)

func WithTaskLogger(logger *zap.Logger) TaskOption {
	return func(opts *TaskOptions) {
		opts.Logger = logger
	}
}

func WithName(name string) TaskOption {
	return func(opts *TaskOptions) {
		opts.Name = name
	}
}

func WithCookie(cookie string) TaskOption {
	return func(opts *TaskOptions) {
		opts.Cookie = cookie
	}
}

func WithWaitTime(waitTime int64) TaskOption {
	return func(opts *TaskOptions) {
		opts.WaitTime = waitTime
	}
}

func WithReload(reload bool) TaskOption {
	return func(opts *TaskOptions) {
		opts.Reload = reload
	}
}

func WithFetcher(f Fetcher) TaskOption {
	return func(opts *TaskOptions) {
		opts.Fetcher = f
	}
}

func WithStorage(s collector.Storage) TaskOption {
	return func(opts *TaskOptions) {
		opts.Storage = s
	}
}

func WithMaxDepth(maxDepth int64) TaskOption {
	return func(opts *TaskOptions) {
		opts.MaxDepth = maxDepth
	}
}
