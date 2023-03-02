package worker

import (
	"crawler/internal/pkg/collector"
	"crawler/internal/pkg/config"
	"crawler/internal/pkg/engine"
	"crawler/internal/pkg/grpc"
	"crawler/internal/pkg/spider"
	"crawler/internal/pkg/store/mysql"
	"crawler/pkg/log"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/spf13/cobra"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
)

var ServiceName string = "go.micro.server.worker"

var WorkerCmd = &cobra.Command{
	Use:   "worker",
	Short: "run worker service.",
	Long:  "run worker service.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Run()
	},
}

var workerID string
var HTTPListenAddress string
var GRPCListenAddress string

func init() {
	WorkerCmd.Flags().StringVar(
		&workerID, "id", "1", "set worker id")
	WorkerCmd.Flags().StringVar(
		&HTTPListenAddress, "http", ":8080", "set worker HTTP listen address")

	WorkerCmd.Flags().StringVar(
		&GRPCListenAddress, "grpc", ":9090", "set worker GRPC listen address")
}

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
	workerConfig := config.GetWorkerConfig(workerID, HTTPListenAddress, GRPCListenAddress)
	logger.Sugar().Infof("grpc server config,%+v", workerConfig)
	reg := etcd.NewRegistry(registry.Addrs(workerConfig.RegistryAddress))

	// start http proxy to GRPC
	go grpc.RunHTTPServer(logger, workerConfig)

	// start grpc server
	grpc.RunGRPCServer(logger, reg, workerConfig)

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
