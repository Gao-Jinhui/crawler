package worker

import (
	"crawler/internal/pkg/collector"
	"crawler/internal/pkg/config"
	"crawler/internal/pkg/engine"
	"crawler/internal/pkg/spider"
	"crawler/internal/pkg/store/mysql"
	"crawler/pkg/log"
	"go.uber.org/zap"
)

func Run() {
	logger := log.NewZapLogger()

	//proxy
	//proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8889"}
	//p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	//if err != nil {
	//	logger.Error("RoundRobinProxySwitcher failed")
	//}

	if err := config.InitConfig(); err != nil {
		logger.Error("failed to init config", zap.Error(err))
	}

	var f spider.Fetcher = spider.NewBrowserFetch(
		spider.WithTimeout(config.GetFetcherTimeout()),
		//collect.WithProxy(p),
	)

	var storage collector.Storage = mysql.NewSqlClient(config.GetMysqlConfig(), mysql.WithLogger(logger))
	if storage == nil {
		return
	}
	taskConfigs := config.GetTaskConfigs()

	seeds := spider.ParseTaskConfigs(logger, f, storage, taskConfigs)

	c := engine.NewCrawler(
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithWorkCount(5),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
	)
	c.Run()
}
