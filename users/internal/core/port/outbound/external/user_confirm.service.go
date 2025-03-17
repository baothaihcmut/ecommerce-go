package external

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/entities"
)

type SendEmailArg struct {
	Email     string
	Code      string
	FirstName string
	LastName  string
}

type UserConfirmService interface {
	StoreUserInfo(context.Context, *entities.User) (string, error)
	GenerateUrlForConfirm(context.Context, string) (string, error)
	IsUserPendingConfirmSignUp(context.Context, string) (bool, error)
	GetUserInfo(context.Context, string) (*entities.User, error)
}
