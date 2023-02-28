package doubanbook

import (
	"crawler/internal/pkg/spider"
)

var DoubanBookRuleTree = &spider.RuleTree{
	Name: "douban_book_list",
	Root: func() ([]*spider.Request, error) {
		roots := []*spider.Request{
			&spider.Request{
				Priority: 1,
				Url:      "https://book.douban.com",
				Method:   "GET",
				RuleName: "数据tag",
			},
		}
		return roots, nil
	},
	Trunk: map[string]*spider.Rule{
		"数据tag": {ParseFunc: ParseTag},
		"书籍列表":  {ParseFunc: ParseBookList},
		"书籍简介":  {ParseFunc: ParseBookDetail},
	},
}
