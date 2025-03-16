package server

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/queue"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/config"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/controllers"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/mailer"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/services"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/gomail.v2"
)

type Server struct {
	rabbitmq *amqp.Connection
	mail     *gomail.Dialer
	cfg 	*config.CoreConfig
}

func NewServer(
	rabbitmq *amqp.Connection,
	mailer *gomail.Dialer,
	cfg *config.CoreConfig,
) *Server {
	return &Server{
		rabbitmq: rabbitmq,
		mail:     mailer,
		cfg: 	cfg,
	}
}
func (s *Server) Start() {
	//init mailer 
	mailer := mailer.NewMailer(s.mail,s.cfg.Mailer)

	authMailService := services.NewAuthMailService(mailer)
	authMailController := controllers.NewAuthMailController(authMailService)
	//queueService
	queueService, err := queue.NewRabbitMqService(s.rabbitmq)
	if err!= nil {
		return
	}
	defer queueService.Close()
	//register queue
	mapCh := make(map[string] chan *amqp.Delivery)
	authMailController.Register(queueService,mapCh)


	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	errCh := make(chan error,100)
	wg.Add(1)
	go func() {
		defer wg.Done()
		authMailController.Run(errCh)
	}()	
	for q,ch := range mapCh{
		wg.Add(1)
		go func ()  {
			defer wg.Done()
			msgs,err:= queueService.Consume(q,"email-service",true)	
			if err!= nil{
				errCh<-err
				return
			}
			for msg := range msgs {
				ch<-&msg
			}
		}()
	}
	wg.Add(1)
	go func ()  {
		wg.Done()
		for err := range errCh {
			fmt.Println(err)
		}	
	}()
	fmt.Println("Server is running")
	<-sigChan
	//close all channel
	for _,ch:= range mapCh{
		close(ch)
	}
	//wait all worker proccess message
	wg.Wait()
	//close error channel
	close(errCh)

	fmt.Println("Graceful shutdown completed.")
}
