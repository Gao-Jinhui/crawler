package master

import (
	"crawler/internal/pkg/config"
	"crawler/internal/pkg/grpc"
	"crawler/pkg/log"
	"go.uber.org/zap"
)

func Run() {
	logger := log.NewZapLogger()
	if err := config.InitConfig(); err != nil {
		logger.Error("failed to init config", zap.Error(err))
	}
	masterConfig := config.GetMasterConfig()
	logger.Sugar().Infof("grpc server config,%+v", masterConfig)
	
	// start http proxy to GRPC
	go grpc.RunHTTPServer(logger, masterConfig)

	// start grpc server
	grpc.RunGRPCServer(logger, masterConfig)

}
