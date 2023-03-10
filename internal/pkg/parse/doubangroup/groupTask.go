package doubangroup

//
//import (
//	"crawler/internal/pkg/spider"
//	"crawler/internal/pkg/spider/request"
//	"crawler/internal/pkg/spider/task"
//	"fmt"
//	"regexp"
//)
//
//var DoubangroupTask = &task.Task{
//	Name:     "find_douban_sun_room",
//	WaitTime: 2,
//	MaxDepth: 5,
//	Cookie:   "douban-fav-remind=1; ll=\"108289\"; bid=kKBun9tYW6s; gr_user_id=a0c87ec1-cbfa-4c20-905b-19af38bae496; viewed=\"5333562_35871233_35519282_30329536_6709783\"; push_noty_num=0; push_doumail_num=0; __utmv=30149280.12137; ct=y; __utmz=30149280.1676992770.30.17.utmcsr=bing|utmccn=(organic)|utmcmd=organic|utmctr=(not provided); __utmc=30149280; dbcl2=\"121370564:YwrHjptOBhc\"; ck=B8yF; frodotk_db=\"0c722bcd6b8b2b5f2556836842a1821d\"; _pk_ref.100001.8cb4=[\"\",\"\",1677053087,\"https://time.geekbang.org/column/article/612328?screen=full\"]; _pk_ses.100001.8cb4=*; __utma=30149280.362716562.1637497006.1677050699.1677053088.33; __utmt=1; _pk_id.100001.8cb4=d57d51bccc0cbbb9.1637496995.25.1677053976.1677051205.; __utmb=30149280.35.5.1677053976398",
//	Rule: task2.RuleTree{
//		Root: func() ([]*request.Request, error) {
//			var roots []*request.Request
//			for i := 0; i < 125; i += 25 {
//				str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
//				roots = append(roots, &request.Request{
//					Priority: 1,
//					Url:      str,
//					Method:   "GET",
//					RuleName: "解析网站URL",
//				})
//			}
//			return roots, nil
//		},
//		Trunk: map[string]*task2.Rule{
//			"解析网站URL": &task2.Rule{ParseFunc: ParseURL},
//			"解析阳台房":   &task2.Rule{ParseFunc: GetSunRoom},
//		},
//	},
//}
//
//func ParseURL(ctx *task2.Context) (task2.ParseResult, error) {
//	re := regexp.MustCompile(urlListRe)
//
//	matches := re.FindAllSubmatch(ctx.Body, -1)
//	result := task2.ParseResult{}
//
//	for _, m := range matches {
//		u := string(m[1])
//		result.Requests = append(
//			result.Requests, &request.Request{
//				Method:   "GET",
//				Task:     ctx.Req.Task,
//				Url:      u,
//				Depth:    ctx.Req.Depth + 1,
//				RuleName: "解析阳台房",
//			})
//	}
//	return result, nil
//}
//
//func GetSunRoom(ctx *task2.Context) (task2.ParseResult, error) {
//	re := regexp.MustCompile(ContentRe)
//
//	ok := re.Match(ctx.Body)
//	if !ok {
//		return task2.ParseResult{
//			Items: []interface{}{},
//		}, nil
//	}
//	result := task2.ParseResult{
//		Items: []interface{}{ctx.Req.Url},
//	}
//	return result, nil
//}
