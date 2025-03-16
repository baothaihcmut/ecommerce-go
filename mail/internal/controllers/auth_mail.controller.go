package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/events"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/queue"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/services"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AuthMailController interface {
	Register(queue.QueueService,map[string]chan *amqp.Delivery) error
	Run(errCh chan error)
}

type AuthMailControllerImpl struct {
	authMailService services.AuthMailService
	msgChs          map[string]chan *amqp.Delivery
}

// Register implements AuthMailController.
func (a *AuthMailControllerImpl) Register(r queue.QueueService,msgChs map[string]chan *amqp.Delivery) error {
	//register event
	if err:= r.CreateQueue("mail-user-signup",true,false) ; err!= nil{
		return err
	}
	if err := r.BindQueue("mail-user-signup","user.signup","user-events"); err != nil{
		return err
	}
	//create chanel for event

	a.msgChs["mail-user-signup"] = make(chan *amqp.Delivery)
	msgChs["mail-user-signup"] = a.msgChs["mail-user-signup"]
	return nil
}

func (a *AuthMailControllerImpl) Run(errCh chan error) {
	localWg := sync.WaitGroup{}
	for q,ch:= range a.msgChs{
		localWg.Add(1)
		go func() {
			defer localWg.Done()
			for msg := range ch{
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				func() {
					defer cancel()
					var err error
					switch q{
					case "mail-user-signup":
						err= a.handleUserSignUpEvent(ctx,msg.Body)
					default:
						err =errors.New("handler not found in Auth service")
					}
					if err!= nil{
						errCh <- err
					}
				}()
			}
		}()
	}
	localWg.Wait()
}

func(a *AuthMailControllerImpl) handleUserSignUpEvent(ctx context.Context, message []byte)error {
	var e events.UserSignUpEvent
	if err := json.Unmarshal(message,&e) ; err!= nil {
		return err
	}
	if err := a.authMailService.SendMailConfirmSignUp(ctx,&e); err!= nil{
		return err
	}
	return nil
}
func NewAuthMailController(svc services.AuthMailService) AuthMailController {
	return &AuthMailControllerImpl{
		authMailService: svc,
		msgChs: make(map[string]chan *amqp.Delivery,100),
	}
}
