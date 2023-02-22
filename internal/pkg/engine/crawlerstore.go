package engine

import (
	"crawler/internal/pkg/collect"
	"crawler/internal/pkg/parse/doubangroup"
)

type CrawlerStore struct {
	list []*collect.Task
	hash map[string]*collect.Task
}

func (c *CrawlerStore) Add(task *collect.Task) {
	c.hash[task.Name] = task
	c.list = append(c.list, task)
}

func init() {
	Store.Add(doubangroup.DoubangroupTask)
	//Store.AddJSTask(doubangroup.DoubangroupJSTask)
}

var Store = &CrawlerStore{
	list: []*collect.Task{},
	hash: map[string]*collect.Task{},
}
