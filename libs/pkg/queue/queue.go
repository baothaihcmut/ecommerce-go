package queue

import amqp "github.com/rabbitmq/amqp091-go"

type QueueService interface {
	Close() error
	CreateQueue(name string, durable, autodelete bool) error
	BindQueue(queueName, binding, exchange string) error
	Consume(queueName, consumer string, autoAck bool) (<-chan amqp.Delivery, error)
}
