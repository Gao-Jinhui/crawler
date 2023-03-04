package engine

import (
	"context"
	"crawler/internal/pkg/collector"
	"crawler/internal/pkg/master"
	"crawler/internal/pkg/spider"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"strings"
	"sync"
)

type Crawler struct {
	id          string
	out         chan spider.ParseResult
	Visited     map[string]bool
	VisitedLock sync.Mutex
	failures    map[string]*spider.Request // 失败请求id -> 失败请求
	failureLock sync.Mutex
	resources   map[string]*master.ResourceSpec
	rlock       sync.Mutex
	etcdCli     *clientv3.Client
	options
}

func NewCrawler(opts ...Option) (*Crawler, error) {
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

	// todo
	// 任务加上默认的采集器与存储器
	//for _, task := range Store.list {
	//	task.Fetcher = c.Fetcher
	//	task.Storage = c.Storage
	//}
	endpoints := []string{c.registryURL}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		return nil, err
	}
	c.etcdCli = cli
	Store.AddTasks(c.TaskSeeds...)
	return c, nil
}

func (c *Crawler) Run(id string, cluster bool) {
	c.id = id
	if !cluster {
		c.handleSeeds()
	}
	go c.loadResource()
	go c.watchResource()
	go c.Schedule()
	for i := 0; i < c.WorkCount; i++ {
		go c.CreateWork()
	}
	c.HandleResult()
}

func (c *Crawler) Schedule() {
	go c.scheduler.Schedule()
}

func (c *Crawler) handleSeeds() {
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

func (c *Crawler) watchResource() {
	watch := c.etcdCli.Watch(context.Background(), master.RESOURCEPATH, clientv3.WithPrefix())
	for w := range watch {
		if w.Err() != nil {
			c.Logger.Error("watch resource failed", zap.Error(w.Err()))
			continue
		}
		if w.Canceled {
			c.Logger.Error("watch resource canceled")
			return
		}
		for _, ev := range w.Events {
			spec, err := master.Decode(ev.Kv.Value)
			if err != nil {
				c.Logger.Error("decode etcd value failed", zap.Error(err))
			}

			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					c.Logger.Info("receive create resource", zap.Any("spec", spec))

				} else if ev.IsModify() {
					c.Logger.Info("receive update resource", zap.Any("spec", spec))
				}
				c.runTasks(spec.Name)
			case clientv3.EventTypeDelete:
				c.Logger.Info("receive delete resource", zap.Any("spec", spec))
			}
		}
	}
}

func getID(assignedNode string) string {
	s := strings.Split(assignedNode, "|")
	if len(s) < 2 {
		return ""
	}
	return s[0]
}

func (c *Crawler) loadResource() error {
	resp, err := c.etcdCli.Get(context.Background(), master.RESOURCEPATH, clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return fmt.Errorf("etcd get failed")
	}

	resources := make(map[string]*master.ResourceSpec)
	for _, kv := range resp.Kvs {
		r, err := master.Decode(kv.Value)
		if err == nil && r != nil {
			id := getID(r.AssignedNode)
			if len(id) > 0 && c.id == id {
				resources[r.Name] = r
			}
		}
	}
	c.Logger.Info("leader init load resource", zap.Int("lenth", len(resources)))
	c.rlock.Lock()
	defer c.rlock.Unlock()
	c.resources = resources
	for _, r := range c.resources {
		c.runTasks(r.Name)
	}

	return nil
}

func (c *Crawler) runTasks(taskName string) {
	task := Store.GetTask(taskName)
	ruletree := Store.GetRuleTree(taskName)
	task.Rule = ruletree //todo
	res, err := ruletree.Root()
	if err != nil {
		c.Logger.Error("get root failed",
			zap.Error(err),
		)
		return
	}
	for _, req := range res {
		req.Task = task
	}
	c.scheduler.Push(res...)
}
