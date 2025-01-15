package config

import (
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/config"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
)

type Config struct {
	Server   ServerConfig        `mapstructure:"server"`
	Database DatabaseConfig      `mapstructure:"db"`
	Consol   ConsulConfig        `mapstructure:"consul"`
	Jwt      JwtConfig           `mapstructure:"jwt"`
	Logger   logger.ConfigLogger `mapstructure:"logger"`
}

type JwtConfig struct {
	AccessTokenSecret  string `mapstructure:"access_token_secret"`
	AccessTokenAge     int    `mapstructure:"access_token_age"`
	RefreshTokenSecret string `mapstructure:"refresh_token_secret"`
	RefreshTokenAge    int    `mapstructure:"refresh_token_age"`
}

type ServerConfig struct {
	Host              string `mapstructure:"host"`
	Port              int    `mapstructure:"port"`
	MaxConnectionIdle int    `mapstructure:"max_connection_idle"`
	MaxConnectionAge  int    `mapstructure:"max_connection_age"`
	Time              int    `mapstructure:"time"`
	Timeout           int    `mapstructure:"time_out"`
}

type DatabaseConfig struct {
	Driver      string `mapstructure:"driver"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Name        string `mapstructure:"name"`
	Ssl         bool   `mapstructure:"ssl"`
	SslMode     string `mapstructure:"ssl_mode"`
	SslCertPath string `mapstructure:"ssl_cert_path"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
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
