package endpoints

import (
	commandHandler "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/handlers"
	queryHandler "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/query/handlers"
)

type UserEnpoints struct {
}

func MakeUserEndpoints(c commandHandler.UserCommandHandler, q queryHandler.UserQueryHandler) UserEnpoints {
	return UserEnpoints{}
}
