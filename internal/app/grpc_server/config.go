package grpc_server

import (
	cfg "go.uber.org/config"
)

type Config struct {
	ListenAddress string `yaml:"listen_address"`
}

func NewServerConfig(provider cfg.Provider) (*Config, error) {
	config := Config{}

	if err := provider.Get("grpc_server").Populate(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
