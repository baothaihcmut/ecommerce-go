package config

import (
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/config"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
)

type Config struct {
	ServerConfig   ServerConfig         `mapstructure:"server"`
	LoggerConfig   logger.ConfigLogger  `mapstructure:"logger"`
	GrpcService    GrpcServiceConfig    `mapstructure:"grpc_service"`
	GrpcConnection GrpcConnectionConfig `mapstructure:"grpc_connection"`
	JaegerConfig   JaegerConfig         `mapstructure:"jaeger"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type GrpcConnectionConfig struct {
	MinConnection         int `mapstructure:"min_connection"`
	MaxConntecion         int `mapstructure:"max_connection"`
	ConnectionIdleTimeOut int `mapstructure:"connection_time_out"`
}

type GrpcServiceConfig struct {
	UserService    string `mapstructure:"user_service"`
	ProductService string `mapstructure:"product_service"`
}
type JaegerConfig struct {
	Endpoint string `mapstructure:"endpoint"`
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
