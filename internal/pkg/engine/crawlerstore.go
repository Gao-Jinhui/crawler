package engine

import (
	"crawler/internal/pkg/parse/doubanbook"
	"crawler/internal/pkg/spider"
)

func init() {
	Store.Add(doubanbook.DoubanBookRuleTree)
	//Store.AddJSTask(doubangroup.DoubangroupJSTask)
}

type CrawlerStore struct {
	ruleList []*spider.RuleTree
	hash     map[string]*spider.RuleTree
}

func (c *CrawlerStore) Add(tree *spider.RuleTree) {
	c.hash[tree.Name] = tree
	c.ruleList = append(c.ruleList, tree)
}

func (c *CrawlerStore) GetRuleTree(name string) *spider.RuleTree {
	return c.hash[name]
}

var Store = &CrawlerStore{
	ruleList: []*spider.RuleTree{},
	hash:     map[string]*spider.RuleTree{},
}
