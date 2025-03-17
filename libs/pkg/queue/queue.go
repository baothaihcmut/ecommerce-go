package queue

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueService interface {
	Close() error
	CreateQueue(name string, durable, autodelete bool) error
	BindQueue(queueName, binding, exchange string) error
	Consume(queueName, consumer string, autoAck bool) (<-chan amqp.Delivery, error)
	Send(context.Context, string, string, any, map[string]string) error
	InitExchange(exchange, exchangeType string) error
}
