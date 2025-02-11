package endpoints

import (
	commandHandler "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/handlers"
)

type UserEnpoints struct {
}

func MakeUserEndpoints(c commandHandler.UserCommandHandler) UserEnpoints {
	return UserEnpoints{}
}
