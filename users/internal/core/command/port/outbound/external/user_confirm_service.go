package external

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/events"
)

type UserConfirmService interface {
	StoreUserInfo(context.Context, *user.User) (string, error)
	PublishSignUpEvent(context.Context, *events.UserSignUpEvent) error
}
