package engine

import (
	"crawler/internal/pkg/collector"
	"crawler/internal/pkg/spider"
	"go.uber.org/zap"
	"sync"
)

type Crawler struct {
	out         chan spider.ParseResult
	Visited     map[string]bool
	VisitedLock sync.Mutex
	failures    map[string]*spider.Request // 失败请求id -> 失败请求
	failureLock sync.Mutex
	options
}

func NewCrawler(opts ...Option) *Crawler {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	c := &Crawler{}
	out := make(chan spider.ParseResult)
	c.Visited = make(map[string]bool, 100)
	c.failures = make(map[string]*spider.Request)
	c.out = out
	c.options = options
	return c
}

func (c *Crawler) Run() {
	go c.Schedule()
	for i := 0; i < c.WorkCount; i++ {
		go c.CreateWork()
	}
	c.HandleResult()
}

func (c *Crawler) Schedule() {
	var reqs []*spider.Request
	for _, taskSeed := range c.TaskSeeds {
		//todo
		if taskSeed == nil {
			continue
		}
		taskSeed.Rule = Store.GetRuleTree(taskSeed.Name)
		rootreqs, err := taskSeed.Rule.Root()
		if err != nil {
			c.Logger.Error("get root failed",
				zap.Error(err),
			)
			continue
		}
		for _, req := range rootreqs {
			req.Task = taskSeed
		}
		reqs = append(reqs, rootreqs...)
	}
	go c.scheduler.Schedule()
	go c.scheduler.Push(reqs...)
}

func (c *Crawler) CreateWork() {
	for {
		req := c.scheduler.Pull()
		if err := req.CheckDepth(); err != nil {
			c.Logger.Error("check failed",
				zap.Error(err),
			)
			continue
		}
		if !req.Task.Reload && c.HasVisited(req) {
			c.Logger.Debug("request has visited",
				zap.String("url:", req.Url),
			)
			continue
		}
		c.StoreVisited(req)

		body, err := req.Fetch()
		if err != nil {
			c.Logger.Error("can't fetch ",
				zap.Error(err),
				zap.String("url", req.Url),
			)
			c.SetFailure(req)
			continue
		}

		if len(body) < 6000 {
			c.Logger.Error("can't fetch ",
				zap.Int("length", len(body)),
				zap.String("url", req.Url),
			)
			c.SetFailure(req)
			continue
		}

		rule := req.Task.Rule.Trunk[req.RuleName]

		res, err := rule.ParseFunc(&spider.Context{
			Body: body,
			Req:  req,
		})

		if err != nil {
			c.Logger.Error("ParseFunc failed ",
				zap.Error(err),
				zap.String("url", req.Url),
			)
			continue
		}
		// put requests into request channel
		if len(res.Requests) > 0 {
			go c.scheduler.Push(res.Requests...)
		}

		c.out <- res
	}
}

func (c *Crawler) HandleResult() {
	for {
		select {
		case res := <-c.out:
			for _, item := range res.Items {
				switch d := item.(type) {
				case *collector.DataCell:
					if err := d.Storage.Save(d); err != nil {
						c.Logger.Error("failed to create", zap.String("error", err.Error()))
					}
				}
				c.Logger.Sugar().Info("get result: ", item)
			}
		}
	}
}

func (c *Crawler) HasVisited(r *spider.Request) bool {
	c.VisitedLock.Lock()
	defer c.VisitedLock.Unlock()
	unique := r.Unique()
	return c.Visited[unique]
}

func (c *Crawler) StoreVisited(reqs ...*spider.Request) {
	c.VisitedLock.Lock()
	defer c.VisitedLock.Unlock()

	for _, r := range reqs {
		unique := r.Unique()
		c.Visited[unique] = true
	}
}

func (c *Crawler) SetFailure(req *spider.Request) {
	if !req.Task.Reload {
		c.VisitedLock.Lock()
		unique := req.Unique()
		delete(c.Visited, unique)
		c.VisitedLock.Unlock()
	}
	c.failureLock.Lock()
	defer c.failureLock.Unlock()
	if _, ok := c.failures[req.Unique()]; !ok {
		// 首次失败时，再重新执行一次
		c.failures[req.Unique()] = req
		c.scheduler.Push(req)
	}
	// todo: 失败2次，加载到失败队列中
}
