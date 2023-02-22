package main

import (
	"crawler/internal/pkg/collect"
	"crawler/internal/pkg/engine"
	"crawler/pkg/log"
	"time"
)

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

	seeds := make([]*collect.Task, 0, 1000)

	seeds = append(seeds, &collect.Task{
		//Name:    "find_douban_sun_room",
		Name:    "douban_book_list",
		Fetcher: f,
	})

	c := engine.NewCrawler(
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithWorkCount(5),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
	)
	c.Run()
}
