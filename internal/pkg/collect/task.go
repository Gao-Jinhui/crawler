package collect

import "time"

type Task struct {
	Url         string
	Cookie      string
	MaxDepth    int
	RootRequest *Request
	WaitTime    time.Duration
	Fetcher     Fetcher
}
