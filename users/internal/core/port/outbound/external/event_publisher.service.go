package external

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/events"
)

type EventPublisherService interface {
	PublishUserSignUpEvent(context.Context, events.UserSignUpEvent) error
}
