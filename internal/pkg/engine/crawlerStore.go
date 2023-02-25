package engine

import (
	"crawler/internal/pkg/collect"
	"crawler/internal/pkg/parse/doubanbook"
	"crawler/internal/pkg/parse/doubangroup"
)

func init() {
	Store.Add(doubangroup.DoubangroupTask)
	Store.Add(doubanbook.DoubanBookTask)
	//Store.AddJSTask(doubangroup.DoubangroupJSTask)
}

type CrawlerStore struct {
	list []*collect.Task
	hash map[string]*collect.Task
}

func (c *CrawlerStore) Add(task *collect.Task) {
	c.hash[task.Name] = task
	c.list = append(c.list, task)
}

var Store = &CrawlerStore{
	list: []*collect.Task{},
	hash: map[string]*collect.Task{},
}
