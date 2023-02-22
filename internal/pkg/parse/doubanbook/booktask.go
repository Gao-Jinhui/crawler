package doubanbook

import (
	"crawler/internal/pkg/collect"
	"time"
)

var DoubanBookTask = &collect.Task{
	Name:     "douban_book_list",
	WaitTime: 1 * time.Second,
	MaxDepth: 5,
	Cookie:   "douban-fav-remind=1; ll=\"108289\"; bid=kKBun9tYW6s; gr_user_id=a0c87ec1-cbfa-4c20-905b-19af38bae496; viewed=\"5333562_35871233_35519282_30329536_6709783\"; push_noty_num=0; push_doumail_num=0; __utmv=30149280.12137; ct=y; __utmz=30149280.1676992770.30.17.utmcsr=bing|utmccn=(organic)|utmcmd=organic|utmctr=(not provided); __utmc=30149280; dbcl2=\"121370564:YwrHjptOBhc\"; ck=B8yF; frodotk_db=\"0c722bcd6b8b2b5f2556836842a1821d\"; _pk_ref.100001.8cb4=[\"\",\"\",1677053087,\"https://time.geekbang.org/column/article/612328?screen=full\"]; _pk_ses.100001.8cb4=*; __utma=30149280.362716562.1637497006.1677050699.1677053088.33; __utmt=1; _pk_id.100001.8cb4=d57d51bccc0cbbb9.1637496995.25.1677053976.1677051205.; __utmb=30149280.35.5.1677053976398",

	Rule: collect.RuleTree{
		Root: func() ([]*collect.Request, error) {
			roots := []*collect.Request{
				&collect.Request{
					Priority: 1,
					Url:      "https://book.douban.com",
					Method:   "GET",
					RuleName: "数据tag",
				},
			}
			return roots, nil
		},
		Trunk: map[string]*collect.Rule{
			"数据tag": &collect.Rule{ParseFunc: ParseTag},
			"书籍列表":  &collect.Rule{ParseFunc: ParseBookList},
			"书籍简介": &collect.Rule{
				ItemFields: []string{
					"书名",
					"作者",
					"页数",
					"出版社",
					"得分",
					"价格",
					"简介",
				},
				ParseFunc: ParseBookDetail,
			},
		},
	},
}
