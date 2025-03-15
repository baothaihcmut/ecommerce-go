package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/entities"
)

type UserRepo interface {
	CreateUser(context.Context, *entities.User) error
	FindUserByEmail(context.Context, string) (*entities.User, error)
	FindUserByPhoneNumber(context.Context, string) (*entities.User, error)
}
