package server

import (
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/queue"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/gomail.v2"
)

type Server struct {
	rabbitmq *amqp.Connection
	mail     *gomail.Dialer
}

func NewServer(
	rabbitmq *amqp.Connection,
	mailer *gomail.Dialer,
) *Server {
	return &Server{
		rabbitmq: rabbitmq,
		mail:     mailer,
	}
}
func (s *Server) Start() {
	//queueService
	queueService, err := queue.NewRabbitMqService()
}
