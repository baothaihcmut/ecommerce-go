package events

import "github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/entities"

type UserSignUpEvent struct {
	User       *entities.User
	ConfrimUrl string
}
