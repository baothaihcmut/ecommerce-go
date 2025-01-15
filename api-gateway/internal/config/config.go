package config

import (
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/config"
)

type Config struct {
	ServerConfig ServerConfig `mapstructure:"server"`
	ConsulConfig ConsulConfig `mapstructure:"consul"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type LoggerConfig struct {
	Mode              string `yaml:"mode" mapstructure:"mode"`
	DisableCaller     bool   `yaml:"disable_caller" mapstructure:"disable_caller"`
	DisableStacktrace bool   `yaml:"disable_stacktrace" mapstructure:"disable_stacktrace"`
	Encoding          string `yaml:"encoding" mapstructure:"encoding"`
	Level             string `yaml:"level" mapstructure:"level"`
	ZapType           string `yaml:"zap_type" mapstructure:"zap_type"`
}

func LoadConfig() *Config {
	cfg, err := config.LoadConfig()
}
