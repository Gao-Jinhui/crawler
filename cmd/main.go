package main

import (
	"crawler/internal/pkg/collect"
	"crawler/internal/pkg/collector"
	"crawler/internal/pkg/config"
	"crawler/internal/pkg/engine"
	"crawler/internal/pkg/store/mysql"
	"crawler/pkg/log"
	"go.uber.org/zap"
	"time"
)

func main() {
	logger := log.NewZapLogger()

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
		collect.WithTimeout(3000*time.Millisecond),
		collect.WithLogger(logger),
		//collect.WithProxy(p),
	)

	var storage collector.Storage = mysql.NewSqlClient(config.GetMysqlConfig(), mysql.WithLogger(logger))
	if storage == nil {
		return
	}
	//book := &model.Book{
	//	Name:      "111111",
	//	Author:    "111111",
	//	Page:      2,
	//	Publisher: "111111",
	//	Score:     "111111",
	//	Price:     "111111",
	//	Intro:     "111111",
	//	Url:       "222222222",
	//}
	//err := storage.Save(&collector.DataCell{Data: map[string]interface{}{
	//	"Data": book,
	//}})

	seeds := make([]*collect.Task, 0, 1000)

	seeds = append(seeds, &collect.Task{
		//Name:    "find_douban_sun_room",
		Name:    "douban_book_list",
		Fetcher: f,
		Storage: storage, //todo
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
