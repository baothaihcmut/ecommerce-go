package config

import (
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/cache"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/logger"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/queue"
)

type CoreConfig struct {
	Server   *ServerConfig   `mapstructure:"server"`
	DB       *DBConfig       `mapstructure:"db"`
	Redis    *cache.RedisConfig    `mapstructure:"redis"`
	RabbitMq *queue.RabbitMqConfig `mapstructure:"rabbitmq"`
	Logger 	 *logger.LoggerConfig `mapstructure:"logger"`
	Jwt *JwtConfig `mapstructure:"jwt"`
}
