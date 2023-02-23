package doubanbook

import (
	"crawler/internal/pkg/collect"
	"time"
)

var DoubanBookTask = &collect.Task{
	Name:     "douban_book_list",
	WaitTime: 2 * time.Second,
	MaxDepth: 5,
	Cookie:   "douban-fav-remind=1; ll=\"108289\"; bid=kKBun9tYW6s; _vwo_uuid_v2=DB23D89A4EBC069790587D76C47D9F426|02ac414d825fe449713e0cb6580492bd; gr_user_id=a0c87ec1-cbfa-4c20-905b-19af38bae496; viewed=\"5333562_35871233_35519282_30329536_6709783\"; push_noty_num=0; push_doumail_num=0; __utmv=30149280.12137; ct=y; dbcl2=\"121370564:YwrHjptOBhc\"; ck=B8yF; ap_v=0,6.0; gr_session_id_22c937bbd8ebd703f2d8e9445f7dfd03=c3bf7c6a-3654-462c-8259-dd0e7c60c9a0; gr_cs1_c3bf7c6a-3654-462c-8259-dd0e7c60c9a0=user_id:1; _pk_ref.100001.3ac3=[\"\",\"\",1677153732,\"https://www.douban.com/\"]; _pk_ses.100001.3ac3=*; __utma=30149280.362716562.1637497006.1677065654.1677153732.36; __utmc=30149280; __utmz=30149280.1677153732.36.18.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmt_douban=1; __utma=81379588.237322096.1670510836.1677065654.1677153732.7; __utmc=81379588; __utmz=81379588.1677153732.7.6.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmt=1; gr_session_id_22c937bbd8ebd703f2d8e9445f7dfd03_c3bf7c6a-3654-462c-8259-dd0e7c60c9a0=true; frodotk_db=\"6827aa64b0f089ad3876ff2e12d19d07\"; _pk_id.100001.3ac3=eafb8e6878deca88.1670510836.7.1677153744.1677065654.; __utmb=30149280.3.10.1677153732; __utmb=81379588.3.10.1677153732",

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
