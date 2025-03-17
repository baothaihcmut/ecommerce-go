package config

type RedisConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}
