package engine

import (
	"crawler/internal/pkg/collect"
	"crawler/internal/pkg/collector"
	"go.uber.org/zap"
	"sync"
)

type Crawler struct {
	out         chan collect.ParseResult
	Visited     map[string]bool
	VisitedLock sync.Mutex
	failures    map[string]*collect.Request // 失败请求id -> 失败请求
	failureLock sync.Mutex
	options
}

func NewCrawler(opts ...Option) *Crawler {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	c := &Crawler{}
	out := make(chan collect.ParseResult)
	c.Visited = make(map[string]bool, 100)
	c.failures = make(map[string]*collect.Request)
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
	var reqs []*collect.Request
	for _, seed := range c.Seeds {
		task := Store.hash[seed.Name]
		task.Fetcher = seed.Fetcher
		task.Storage = seed.Storage
		rootreqs, err := task.Rule.Root()
		if err != nil {
			c.Logger.Error("get root failed",
				zap.Error(err),
			)
			continue
		}
		for _, req := range rootreqs {
			req.Task = task
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
			//fmt.Println(req.Depth)
			//fmt.Println(req.Task.MaxDepth)
			continue
		}
		if !req.Task.Reload && c.HasVisited(req) {
			c.Logger.Debug("request has visited",
				zap.String("url:", req.Url),
			)
			continue
		}
		c.StoreVisited(req)

		body, err := req.Task.Fetcher.Get(req)
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

		result, err := rule.ParseFunc(&collect.Context{
			body,
			req,
		})

		if err != nil {
			c.Logger.Error("ParseFunc failed ",
				zap.Error(err),
				zap.String("url", req.Url),
			)
			continue
		}

		if len(result.Requests) > 0 {
			go c.scheduler.Push(result.Requests...)
		}

		c.out <- result
	}
}

func (c *Crawler) HandleResult() {
	for {
		select {
		case result := <-c.out:
			for _, item := range result.Items {
				switch d := item.(type) {
				case *collector.DataCell:
					name := d.GetTaskName()
					task := Store.hash[name]
					if err := task.Storage.Save(d); err != nil {
						c.Logger.Error("failed to create", zap.String("error", err.Error()))
					}
				}
				c.Logger.Sugar().Info("get result: ", item)
			}
		}
	}
}

func (c *Crawler) HasVisited(r *collect.Request) bool {
	c.VisitedLock.Lock()
	defer c.VisitedLock.Unlock()
	unique := r.Unique()
	return c.Visited[unique]
}

func (c *Crawler) StoreVisited(reqs ...*collect.Request) {
	c.VisitedLock.Lock()
	defer c.VisitedLock.Unlock()

	for _, r := range reqs {
		unique := r.Unique()
		c.Visited[unique] = true
	}
}

func (c *Crawler) SetFailure(req *collect.Request) {
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
