package services

import (
	"context"

	commonEvent "github.com/baothaihcmut/Ecommerce-go/libs/pkg/events"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/queue"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/events"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/outbound/external"
)

type EventPublisherService struct {
	queueService queue.QueueService
}

// PublishUserSignUpEvent implements external.EventPublisherService.
func (e *EventPublisherService) PublishUserSignUpEvent(ctx context.Context, domainEvent events.UserSignUpEvent) error {
	signUpEvent := commonEvent.UserSignUpEvent{
		Id:          domainEvent.User.Id,
		Email:       domainEvent.User.Email,
		FirstName:   domainEvent.User.FirstName,
		LastName:    domainEvent.User.LastName,
		PhoneNumber: domainEvent.User.PhoneNumber,
		ConfirmUrl:  domainEvent.ConfrimUrl,
	}
	err := e.queueService.Send(ctx, "user-events", "user.signup", signUpEvent, nil)
	if err != nil {
		return err
	}
	return nil
}

func NewEventPublisherService(queueService queue.QueueService) external.EventPublisherService {
	return &EventPublisherService{
		queueService: queueService,
	}
}
