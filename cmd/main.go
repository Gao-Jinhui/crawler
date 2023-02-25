package main

import (
	"crawler/internal/pkg/collect"
	"crawler/internal/pkg/collector"
	"crawler/internal/pkg/config"
	"crawler/internal/pkg/engine"
	"crawler/internal/pkg/limiter"
	"crawler/internal/pkg/store/mysql"
	"crawler/pkg/log"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	logger := log.NewZapLogger()

	secondLimiter := rate.NewLimiter(limiter.Per(1, 2*time.Second), 1)
	minuteLimiter := rate.NewLimiter(limiter.Per(20, 1*time.Minute), 20)
	multiLimiter := limiter.MultiLimiter(secondLimiter, minuteLimiter)

	//proxy
	//proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8889"}
	//p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	//if err != nil {
	//	logger.Error("RoundRobinProxySwitcher failed")
	//}
	if err := config.InitConfig(); err != nil {
		logger.Error("failed to init config", zap.String("err", err.Error()))
	}

	var f collect.Fetcher = collect.NewBrowserFetch(
		collect.WithTimeout(3500*time.Millisecond),
		collect.WithLogger(logger),
		//collect.WithProxy(p),
	)

	var storage collector.Storage = mysql.NewSqlClient(config.GetMysqlConfig(), mysql.WithLogger(logger))
	if storage == nil {
		return
	}

	seeds := make([]*collect.Task, 0, 1000)

	seeds = append(seeds, &collect.Task{
		//Name:    "find_douban_sun_room",
		Name:    "douban_book_list",
		Fetcher: f,
		Storage: storage,
		Limit:   multiLimiter,
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
