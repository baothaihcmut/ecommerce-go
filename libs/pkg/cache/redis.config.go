package cache


type RedisConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Username string `mapstructure:"user_name"`
}
