package collect

import "time"

type Task struct {
	//Property
	Name        string
	Url         string
	Cookie      string
	MaxDepth    int
	RootRequest *Request
	WaitTime    time.Duration
	Fetcher     Fetcher
	Reload      bool // 网站是否可以重复爬取
	Rule        RuleTree
}
