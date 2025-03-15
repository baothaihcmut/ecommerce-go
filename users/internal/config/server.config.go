package config

type ServerConfig struct {
	Port              int `mapstructure:"port"`
	MaxConnectionIdle int `mapstructure:"max_connection_idle"`
}
