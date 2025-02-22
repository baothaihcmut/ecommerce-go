package repositories

import (
	"context"
	"database/sql"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
)

type UserRepository interface {
	CreateUser(context.Context, *user.User, *sql.Tx) error
	UpdateUser(context.Context, *user.User, *sql.Tx) error
	ActivateShopOwner(context.Context, *user.User, *sql.Tx) error
	UpdateUserAddress(ctx context.Context, user *user.User, t *sql.Tx) (err error)

	FindById(context.Context, valueobject.UserId) (*user.User, error)
	FindByEmail(context.Context, valueobject.Email) (*user.User, error)
	CheckEmailExist(context.Context, valueobject.Email) (bool, error)
	CheckPhoneNumberExist(context.Context, valueobject.PhoneNumber) (bool, error)
}
