package config

type CoreConfig struct {
	Server   *ServerConfig   `mapstructure:"server"`
	DB       *DBConfig       `mapstructure:"db"`
	Redis    *RedisConfig    `mapstructure:"redis"`
	RabbitMq *RabbitMqConfig `mapstructure:"rabbitmq"`
}
