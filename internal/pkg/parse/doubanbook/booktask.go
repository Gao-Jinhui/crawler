package doubanbook

import (
	"crawler/internal/pkg/collect"
	"time"
)

var DoubanBookTask = &collect.Task{
	Name:     "douban_book_list",
	WaitTime: 1 * time.Second,
	MaxDepth: 5,
	//Cookie:   "bid=-UXUw--yL5g; dbcl2=\"214281202:q0BBm9YC2Yg\"; __yadk_uid=jigAbrEOKiwgbAaLUt0G3yPsvehXcvrs; push_noty_num=0; push_doumail_num=0; __utmz=30149280.1665849857.1.1.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmv=30149280.21428; ck=SAvm; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1665925405%2C%22https%3A%2F%2Faccounts.douban.com%2F%22%5D; _pk_ses.100001.8cb4=*; __utma=30149280.2072705865.1665849857.1665849857.1665925407.2; __utmc=30149280; __utmt=1; __utmb=30149280.23.5.1665925419338; _pk_id.100001.8cb4=fc1581490bf2b70c.1665849856.2.1665925421.1665849856.",
	Cookie: "douban-fav-remind=1; ll=\"108289\"; bid=kKBun9tYW6s; _vwo_uuid_v2=DB23D89A4EBC069790587D76C47D9F426|02ac414d825fe449713e0cb6580492bd; gr_user_id=a0c87ec1-cbfa-4c20-905b-19af38bae496; viewed=\"5333562_35871233_35519282_30329536_6709783\"; push_noty_num=0; push_doumail_num=0; __utmv=30149280.12137; ct=y; dbcl2=\"121370564:zG7W3OjU/8M\"; ck=Xir_; __utmc=30149280; __utmc=81379588; frodotk_db=\"3b17d6504bff81eab14782f0e3fe85be\"; gr_session_id_22c937bbd8ebd703f2d8e9445f7dfd03=af1331d4-5073-46d4-8830-b46fe6fe4b54; gr_cs1_af1331d4-5073-46d4-8830-b46fe6fe4b54=user_id:1; _pk_ref.100001.3ac3=[\"\",\"\",1677225258,\"https://www.bing.com/\"]; _pk_id.100001.3ac3=eafb8e6878deca88.1670510836.10.1677225258.1677211615.; _pk_ses.100001.3ac3=*; __utma=30149280.362716562.1637497006.1677211615.1677225258.39; __utmz=30149280.1677225258.39.20.utmcsr=bing|utmccn=(organic)|utmcmd=organic|utmctr=(not provided); __utmt_douban=1; __utmb=30149280.1.10.1677225258; __utma=81379588.237322096.1670510836.1677211615.1677225258.11; __utmz=81379588.1677225258.11.8.utmcsr=bing|utmccn=(organic)|utmcmd=organic|utmctr=(not provided); __utmt=1; __utmb=81379588.1.10.1677225258; ap_v=0,6.0; gr_session_id_22c937bbd8ebd703f2d8e9445f7dfd03_af1331d4-5073-46d4-8830-b46fe6fe4b54=true",

	Rule: collect.RuleTree{
		Root: func() ([]*collect.Request, error) {
			roots := []*collect.Request{
				&collect.Request{
					Priority: 1,
					Url:      "https://book.douban.com",
					Method:   "GET",
					RuleName: "数据tag",
				},
				//&collect.Request{
				//	Priority: 1,
				//	Url:      "https://book.douban.com/tag/%E5%B0%8F%E8%AF%B4",
				//	Method:   "GET",
				//	RuleName: "书籍列表",
				//},
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
