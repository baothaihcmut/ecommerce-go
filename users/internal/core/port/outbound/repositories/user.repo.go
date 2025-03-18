package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/entities"
	"github.com/jackc/pgx/v5"
)

type UserRepo interface {
	CreateUser(context.Context, *entities.User) error
	FindUserByEmail(context.Context, string) (*entities.User, error)
	FindUserByPhoneNumber(context.Context, string) (*entities.User, error)
	WithTx(pgx.Tx) UserRepo
}
