package config

import "github.com/baothaihcmut/Ecommerce-go/libs/pkg/logger"

type CoreConfig struct {
	Server   *ServerConfig   `mapstructure:"server"`
	DB       *DBConfig       `mapstructure:"db"`
	Redis    *RedisConfig    `mapstructure:"redis"`
	RabbitMq *RabbitMqConfig `mapstructure:"rabbitmq"`
	Logger 	 *logger.LoggerConfig `mapstructure:"logger"`
}
