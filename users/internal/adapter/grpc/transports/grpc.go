package transports

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/endpoints"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/request"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/response"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	gt "github.com/go-kit/kit/transport/grpc"
)

type GrpcServer struct {
	proto.UnimplementedUserServiceServer
	createUser gt.Handler
}

// CreateUser implements proto.UserServiceServer.
func (g *GrpcServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	_, resp, err := g.createUser.ServeGRPC(ctx, req)
}

// mustEmbedUnimplementedUserServiceServer implements proto.UserServiceServer.

func NewGrpcServer(
	endpoints endpoints.UserEnpoints,
	userRequestMapper request.UserRequestMapper,
	userResponseMapper response.UserResponseMapper,
) proto.UserServiceServer {
	return &GrpcServer{
		createUser: gt.NewServer(
			endpoints.CreateUser,
			userRequestMapper.ToCreateUserCommand,
			userResponseMapper.ToCreateUserResponse,
		),
	}
}
