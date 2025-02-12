package config

import (
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/config"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
)

type Config struct {
	Server ServerConfig        `mapstructure:"server"`
	Mongo  MongoConfig         `mapstructure:"mongo"`
	Jaeger JaegerConfig        `mapstructure:"jaeger"`
	S3     S3Config            `mapstructure:"s3"`
	Logger logger.ConfigLogger `mapstructure:"logger"`
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
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

type JaegerConfig struct {
	Endpoint string `mapstructure:"endpoint"`
}

type S3Config struct {
	Bucket          string `mapstructure:"bucket"`
	StorageProvider string `mapstructure:"storage_provider"`
	AccessKey       string `mapstructure:"access_key"`
	SecretKey       string `mapstructure:"secret_key"`
	Region          string `mapstructure:"region"`
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
