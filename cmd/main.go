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
	cookie := "douban-fav-remind=1; ll=\"108289\"; bid=kKBun9tYW6s; gr_user_id=a0c87ec1-cbfa-4c20-905b-19af38bae496; viewed=\"5333562_35871233_35519282_30329536_6709783\"; push_noty_num=0; push_doumail_num=0; __utmv=30149280.12137; ct=y; __utmz=30149280.1676992770.30.17.utmcsr=bing|utmccn=(organic)|utmcmd=organic|utmctr=(not provided); _pk_ref.100001.8cb4=[\"\",\"\",1677046350,\"https://time.geekbang.org/column/article/612328?screen=full\"]; _pk_ses.100001.8cb4=*; ap_v=0,6.0; __utma=30149280.362716562.1637497006.1676992770.1677046353.31; __utmc=30149280; __utmt=1; dbcl2=\"121370564:YwrHjptOBhc\"; ck=B8yF; _pk_id.100001.8cb4=d57d51bccc0cbbb9.1637496995.23.1677047481.1676626109.; __utmb=30149280.14.5.1677047481110"
	for i := 0; i <= 100; i += 25 {
		url := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		seeds = append(seeds, &collect.Task{
			Url:      url,
			WaitTime: time.Second,
			MaxDepth: 5,
			Fetcher:  f,
			Cookie:   cookie,
			RootRequest: &collect.Request{
				Priority:  1,
				Method:    "GET",
				ParseFunc: doubangroup.ParseURL,
			},
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
