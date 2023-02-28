package spider

type TaskConfig struct {
	Name     string
	Cookie   string
	WaitTime int64
	Reload   bool
	MaxDepth int64
	Fetcher  string
	Limiters []LimitConfig
}

type LimitConfig struct {
	EventCount int
	EventDur   int // 秒
	Bucket     int // 桶大小
}
