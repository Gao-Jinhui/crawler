package grpc

import (
	"crawler/internal/pkg/config"
	"crawler/internal/pkg/master"
	proto "crawler/internal/pkg/proto/crawler"
	grpccli "github.com/go-micro/plugins/v4/client/grpc"
	"github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go.uber.org/zap"
	"time"
)

func RunMasterGRPCServer(m *master.Master, logger *zap.Logger, reg registry.Registry, cfg config.ServerConfig) {
	service := micro.NewService(
		micro.Server(grpc.NewServer(
			server.Id(cfg.ID),
		)),
		micro.Address(cfg.GRPCListenAddress),
		micro.Registry(reg),
		micro.RegisterTTL(time.Duration(cfg.RegisterTTL)*time.Second),
		micro.RegisterInterval(time.Duration(cfg.RegisterInterval)*time.Second),
		micro.WrapHandler(logWrapper(logger)),
		micro.Name(cfg.Name),
		micro.Client(grpccli.NewClient()),
	)

	cl := proto.NewCrawlerMasterService(cfg.Name, service.Client())
	m.SetForwardCli(cl)

	// 设置micro 客户端默认超时时间为10秒钟
	if err := service.Client().Init(client.RequestTimeout(time.Duration(cfg.ClientTimeOut) * time.Second)); err != nil {
		logger.Sugar().Error("micro client init error. ", zap.String("error:", err.Error()))

		return
	}

	service.Init()

	if err := proto.RegisterCrawlerMasterHandler(service.Server(), m); err != nil {
		logger.Fatal("register handler failed", zap.Error(err))
	}

	if err := service.Run(); err != nil {
		logger.Fatal("grpc server stop", zap.Error(err))
	}
}
