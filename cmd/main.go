package main

import (
	"crawler/internal/pkg/collect"
	"crawler/internal/pkg/engine"
	"crawler/internal/pkg/parse/doubangroup"
	"crawler/pkg/log"
	"fmt"
	"time"
)

//tag v0.0.9
func main() {
	logger := log.GetLogger()

	// proxy
	//proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8889"}
	//p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	//if err != nil {
	//	logger.Error("RoundRobinProxySwitcher failed")
	//}

	var f collect.Fetcher = collect.NewBrowserFetch(
		collect.WithTimeout(3000*time.Millisecond),
		collect.WithLogger(logger),
	)

	var seeds []*collect.Task
	cookie := "douban-fav-remind=1; ll=\"108289\"; bid=kKBun9tYW6s; gr_user_id=a0c87ec1-cbfa-4c20-905b-19af38bae496; viewed=\"5333562_35871233_35519282_30329536_6709783\"; ap_v=0,6.0; __utmc=30149280; __utmz=30149280.1676362309.15.14.utmcsr=time.geekbang.org|utmccn=(referral)|utmcmd=referral|utmcct=/column/article/612328; dbcl2=\"121370564:YwrHjptOBhc\"; ck=B8yF; push_noty_num=0; push_doumail_num=0; __utmv=30149280.12137; _pk_ref.100001.8cb4=[\"\",\"\",1676367627,\"https://time.geekbang.org/column/article/612328\"]; _pk_ses.100001.8cb4=*; __utma=30149280.362716562.1637497006.1676362309.1676367628.16; _pk_id.100001.8cb4=d57d51bccc0cbbb9.1637496995.10.1676368492.1676364189.; __utmt=1; __utmb=30149280.12.6.1676368493194"
	for i := 0; i <= 100; i += 25 {
		url := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		seeds = append(seeds, &collect.Task{
			Url:         url,
			WaitTime:    time.Second,
			MaxDepth:    5,
			Fetcher:     f,
			Cookie:      cookie,
			RootRequest: &collect.Request{ParseFunc: doubangroup.ParseURL},
		})
	}

	c := engine.NewCrawler(
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithWorkCount(5),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
	)
	c.Run()
}
