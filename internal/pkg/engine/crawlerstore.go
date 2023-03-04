package engine

import (
	"crawler/internal/pkg/parse/doubanbook"
	"crawler/internal/pkg/spider"
)

func init() {
	Store.AddRuleTrees(doubanbook.DoubanBookRuleTree)
	//Store.AddJSTask(doubangroup.DoubangroupJSTask)
}

type CrawlerStore struct {
	ruleTreeHash map[string]*spider.RuleTree
	taskHash     map[string]*spider.Task
}

func (c *CrawlerStore) AddRuleTrees(trees ...*spider.RuleTree) {
	for _, tree := range trees {
		c.ruleTreeHash[tree.Name] = tree
	}
}

func (c *CrawlerStore) AddTasks(tasks ...*spider.Task) {
	for _, task := range tasks {
		if task != nil {
			c.taskHash[task.Name] = task
		}
	}
}

func (c *CrawlerStore) GetRuleTree(name string) *spider.RuleTree {
	return c.ruleTreeHash[name]
}

func (c *CrawlerStore) GetTask(name string) *spider.Task {
	return c.taskHash[name]
}

var Store = &CrawlerStore{
	taskHash:     map[string]*spider.Task{},
	ruleTreeHash: map[string]*spider.RuleTree{},
}
