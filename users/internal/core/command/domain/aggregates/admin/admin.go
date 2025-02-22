package admin

import (
	"time"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
	"github.com/google/uuid"
)

type Admin struct {
	Id                  valueobject.UserId
	Email               valueobject.Email
	Password            valueobject.Password
	PhoneNumber         valueobject.PhoneNumber
	CurrentRefreshToken string
	FirstName           string
	LastName            string
	LastLoginTime       time.Time
}

func NewAdmin(
	email valueobject.Email,
	password valueobject.Password,
	phoneNumber valueobject.PhoneNumber,
	FirstName string,
	LastName string,
) (*Admin, error) {
	id, err := valueobject.NewUserId(uuid.New())
	if err != nil {
		return nil, err
	}
	return &Admin{
		Id:          *id,
		Email:       email,
		Password:    password,
		PhoneNumber: phoneNumber,
		FirstName:   FirstName,
		LastName:    LastName,
	}, nil
}

func (a *Admin) LogIn(password string) error {
	if !a.Password.Compare(password) {
		return user.ErrBadCredencial
	}
	a.LastLoginTime = time.Now()
	return nil
}
func (a *Admin) SetCurrentRefreshToken(token string) {
	a.CurrentRefreshToken = token
}
