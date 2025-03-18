package initialize

import (
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/queue"
	"github.com/rabbitmq/amqp091-go"
)

func InitializeRabbitMq(cfg *queue.RabbitMqConfig) (*amqp091.Connection, error) {
	return queue.ConnectRabbitMq(cfg.Username, cfg.Password, cfg.Endpoint, cfg.Vhost,cfg.IsSecure)
}
