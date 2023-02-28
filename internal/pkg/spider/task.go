package spider

import (
	"crawler/internal/pkg/collector"
	"crawler/internal/pkg/limiter"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"time"
)

type Task struct {
	Rule *RuleTree
	TaskOptions
}

func NewTask(opts ...TaskOption) *Task {
	options := DefaultTaskOptions
	for _, opt := range opts {
		opt(&options)
	}
	d := &Task{}
	d.TaskOptions = options
	return d
}
func ParseTaskConfigs(logger *zap.Logger, f Fetcher, s collector.Storage, cfgs []TaskConfig) []*Task {
	tasks := make([]*Task, 1000)
	for _, cfg := range cfgs {
		task := NewTask(
			WithName(cfg.Name),
			WithReload(cfg.Reload),
			WithMaxDepth(cfg.MaxDepth),
			WithCookie(cfg.Cookie),
			WithTaskLogger(logger),
			WithStorage(s),
			WithWaitTime(cfg.WaitTime),
		)
		var limits []limiter.RateLimiter
		if len(cfg.Limiters) > 0 {
			for _, lcfg := range cfg.Limiters {
				// speed limiter
				l := rate.NewLimiter(limiter.Per(lcfg.EventCount, time.Duration(lcfg.EventDur)*time.Second), lcfg.Bucket)
				limits = append(limits, l)
			}
			multiLimiter := limiter.MultiLimiter(limits...)
			task.Limit = multiLimiter
		}
		switch cfg.Fetcher {
		case "browser":
			task.Fetcher = f
		}
		tasks = append(tasks, task)
	}
	return tasks
}
