package config

import (
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/logger"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/queue"
)

type CoreConfig struct {
	Mailer   *MailerConfig   `mapstructure:"mailer"`
	RabbitMq *queue.RabbitMqConfig `mapstructure:"rabbitmq"`
	Logger   *logger.LoggerConfig `mapstructure:"logger"`
}
