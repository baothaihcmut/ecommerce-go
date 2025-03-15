package config

type CoreConfig struct {
	Server *ServerConfig `mapstructure:"server"`
	DB     *DBConfig     `mapstructure:"db"`
}
