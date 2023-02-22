package collect

import "time"

type Property struct {
	Name     string        `json:"name"` // 任务名称，应保证唯一性
	Url      string        `json:"url"`
	Cookie   string        `json:"cookie"`
	WaitTime time.Duration `json:"wait_time"`
	Reload   bool          `json:"reload"` // 网站是否可以重复爬取
	MaxDepth int64         `json:"max_depth"`
}
