package grpc

import (
	"context"
	"crawler/internal/pkg/config"
	proto "crawler/internal/pkg/proto/crawler"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func RunMasterHTTPServer(logger *zap.Logger, cfg config.ServerConfig) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := proto.RegisterCrawlerMasterGwFromEndpoint(ctx, mux, cfg.GRPCListenAddress, opts); err != nil {
		logger.Fatal("Register backend grpc server endpoint failed", zap.Error(err))
	}
	zap.S().Infof("start http server listening on %v proxy to grpc server;%v", cfg.HTTPListenAddress, cfg.GRPCListenAddress)
	if err := http.ListenAndServe(cfg.HTTPListenAddress, mux); err != nil {
		logger.Fatal("http listenAndServe failed", zap.Error(err))
	}
}
