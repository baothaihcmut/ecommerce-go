package initialize

import (
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/queue"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

func InitializeRabbitMq(cfg *config.RabbitMqConfig) (*amqp091.Connection, error) {
	return queue.ConnectRabbitMq(cfg.Username, cfg.Password, cfg.Endpoint, cfg.Vhost)
}
