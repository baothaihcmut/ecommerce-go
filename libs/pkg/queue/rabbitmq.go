package queue

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqServiceImpl struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// Send implements QueueService.
func (r *RabbitMqServiceImpl) Send(ctx context.Context, exchange string, binding string, data any, headers map[string]string) error {
	//encode body
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	table := amqp.Table{}
	if headers != nil {
		for k, v := range headers {
			table[k] = v
		}
	}
	err = r.ch.PublishWithContext(
		ctx, exchange,
		binding,
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			Headers:      table,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func ConnectRabbitMq(username, password, host, vhost string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewRabbitMqService(conn *amqp.Connection) (QueueService, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMqServiceImpl{
		conn: conn,
		ch:   ch,
	}, nil
}

func (r *RabbitMqServiceImpl) Close() error {
	return r.conn.Close()
}

func (rc *RabbitMqServiceImpl) CreateQueue(name string, durable, autodelete bool) error {
	_, err := rc.ch.QueueDeclare(name, durable, autodelete, false, false, nil)
	return err
}

func (rc *RabbitMqServiceImpl) BindQueue(queueName, binding, exchange string) error {
	return rc.ch.QueueBind(queueName, binding, exchange, false, nil)
}
func (rc *RabbitMqServiceImpl) Consume(queueName, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queueName, consumer, autoAck, false, false, false, nil)
}
func (rc *RabbitMqServiceImpl) InitExchange(exchange, exchangeType string) error {
	err := rc.ch.ExchangeDeclare(
		exchange,
		exchangeType,
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}
