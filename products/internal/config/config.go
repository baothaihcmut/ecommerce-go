package config

import (
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/config"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
)

type Config struct {
	Server ServerConfig        `mapstruct:"server"`
	Mongo  MongoConfig         `mapstruct:"mongo"`
	Jaeger JaegerConfig        `mapstruct:"jaeger"`
	Logger logger.ConfigLogger `mapstruct:"logger"`
}

type ServerConfig struct {
	Host              string `mapstructure:"host"`
	Port              int    `mapstructure:"port"`
	MaxConnectionIdle int    `mapstructure:"max_connection_idle"`
	MaxConnectionAge  int    `mapstructure:"max_connection_age"`
	Time              int    `mapstructure:"time"`
	Timeout           int    `mapstructure:"time_out"`
}

type MongoConfig struct {
	URI      string `mapstruct:"uri"`
	Database string `mapstruct:"database"`
}

type JaegerConfig struct {
	Address string `mapstruct:"address"`
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
