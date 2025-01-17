package config

import (
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/config"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
)

type Config struct {
	ServerConfig ServerConfig        `mapstructure:"server"`
	ConsulConfig ConsulConfig        `mapstructure:"consul"`
	LoggerConfig logger.ConfigLogger `mapstructure:"logger"`
	GrpcService  GrpcServiceConfig   `mapstructure:"grpc_service"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type GrpcServiceConfig struct {
	UserService string `mapstructure:"user_service"`
}

func LoadConfig() *Config {
	cfg := &Config{}
	cfgInterface, err := config.LoadConfig(cfg, "./config")
	if err != nil {
		panic(err)
	}
	cfg = cfgInterface.(*Config)
	return cfg
}
