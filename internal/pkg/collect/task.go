package collect

import (
	"crawler/internal/pkg/collector"
	"crawler/internal/pkg/limiter"
)

type Task struct {
	//Property
	Name        string
	Url         string
	Cookie      string
	MaxDepth    int
	RootRequest *Request
	WaitTime    int64
	Fetcher     Fetcher
	Reload      bool // 网站是否可以重复爬取
	Rule        RuleTree
	Storage     collector.Storage
	Limit       limiter.RateLimiter
}
