package config

type WorkerConfig struct {
	GRPCListenAddress string
	HTTPListenAddress string
	ID                string
	RegistryAddress   string
	RegisterTTL       int
	RegisterInterval  int
	Name              string
	ClientTimeOut     int
}
