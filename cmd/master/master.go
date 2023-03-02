package master

import (
	"crawler/internal/pkg/config"
	"crawler/internal/pkg/grpc"
	"crawler/internal/pkg/master"
	"crawler/pkg/log"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/spf13/cobra"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
)

var MasterCmd = &cobra.Command{
	Use:   "master",
	Short: "run master service.",
	Long:  "run master service.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Run()
	},
}

var masterID string
var HTTPListenAddress string
var GRPCListenAddress string

func init() {
	MasterCmd.Flags().StringVar(
		&masterID, "id", "1", "set master id")
	MasterCmd.Flags().StringVar(
		&HTTPListenAddress, "http", ":8081", "set master HTTP listen address")

	MasterCmd.Flags().StringVar(
		&GRPCListenAddress, "grpc", ":9091", "set master GRPC listen address")
}

func Run() {
	logger := log.NewZapLogger()
	if err := config.InitConfig(); err != nil {
		logger.Error("failed to init config", zap.Error(err))
	}
	masterConfig := config.GetMasterConfig(masterID, HTTPListenAddress, GRPCListenAddress)
	logger.Sugar().Infof("grpc server config,%+v", masterConfig)
	reg := etcd.NewRegistry(registry.Addrs(masterConfig.RegistryAddress))
	master.New(
		masterID,
		master.WithLogger(logger.Named("master")),
		master.WithGRPCAddress(GRPCListenAddress),
		master.WithregistryURL(masterConfig.RegistryAddress),
		master.WithRegistry(reg),
	)
	// start http proxy to GRPC
	go grpc.RunHTTPServer(logger, masterConfig)

	// start grpc server
	grpc.RunGRPCServer(logger, reg, masterConfig)

}
