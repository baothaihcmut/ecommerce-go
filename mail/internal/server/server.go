package server

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/logger"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/queue"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/config"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/controllers"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/mailer"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/services"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Server struct {
	rabbitmq *amqp.Connection
	mail     *gomail.Dialer
	logrus 	*logrus.Logger
	cfg      *config.CoreConfig
}

func NewServer(
	rabbitmq *amqp.Connection,
	mailer *gomail.Dialer,
	logrus *logrus.Logger,
	cfg *config.CoreConfig,
) *Server {
	return &Server{
		rabbitmq: rabbitmq,
		mail:     mailer,
		cfg:      cfg,
		logrus: logrus,
	}
}
func (s *Server) Start() {
	loggerService := logger.NewLogger(s.logrus)
	//init mailer
	mailer := mailer.NewMailer(s.mail, s.cfg.Mailer,loggerService)

	authMailService := services.NewAuthMailService(mailer)
	authMailController := controllers.NewAuthMailController(authMailService)
	//queueService
	queueService, err := queue.NewRabbitMqService(s.rabbitmq)
	if err != nil {
		return
	}
	defer queueService.Close()
	//register queue
	mapCh := make(map[string]chan *amqp.Delivery)
	authMailController.Register(queueService, mapCh)

	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	errCh := make(chan error, 100)
	wg.Add(1)
	go func() {
		defer wg.Done()
		authMailController.Run(errCh)
	}()
	for q, ch := range mapCh {
		wg.Add(1)
		go func(q string, ch chan *amqp.Delivery) {
			defer wg.Done()
			msgs, err := queueService.Consume(q, "email-service", true)
			if err != nil {
				errCh <- err
				return
			}
			go func() {
				for msg := range msgs {
					loggerService.Info(context.Background(),map[string]any{
						"queue": q,
						"route_key": msg.RoutingKey,
						"exchange": msg.Exchange,
						"time_stamp": msg.Timestamp,
					},"Incomming event")
					ch <- &msg
				}
			}()
		}(q, ch)
	}
	wg.Add(1)
	go func() {
		wg.Done()
		for err := range errCh {
			loggerService.Errorf(context.Background(),nil,"Error: %v",err)
		}
	}()
	loggerService.Info(context.Background(),nil,"Server is running")
	<-sigChan
	//close all channel
	for _, ch := range mapCh {
		close(ch)
	}
	//wait all worker proccess message
	wg.Wait()
	//close error channel
	close(errCh)

	loggerService.Info(context.Background(),nil,"Graceful shutdown completed.")
}
