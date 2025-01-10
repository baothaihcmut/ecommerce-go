package transports

import (
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/endpoints"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/request"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/response"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
)

type UserServer struct {
}

func NewUserServer(
	endpoints endpoints.UserEnpoints,
	requestMapper request.UserRequestMapper,
	responseMapper response.UserResponseMapper,
) proto.UserServiceServer {
	return &UserServer{}
}
