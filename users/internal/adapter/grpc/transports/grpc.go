package transports

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/endpoints"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/request"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/response"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/utils"
	gt "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
)

type GrpcServer struct {
	proto.UnimplementedUserServiceServer
	createUser gt.Handler
}

// CreateUser implements proto.UserServiceServer.
func (g *GrpcServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	_, resp, err := g.createUser.ServeGRPC(ctx, req)
	if err != nil {
		return &proto.CreateUserResponse{Status: &proto.Status{
			Success: false,
			Code:    utils.MapErrorToGrpcStatus(err).String(),
			Message: err.Error(),
			Details: []string{err.Error()},
		},
			Data: nil,
		}, err
	}
	return &proto.CreateUserResponse{
		Data: resp.(*proto.CreateUserData),
		Status: &proto.Status{
			Success: true,
			Code:    codes.OK.String(),
			Message: "Success",
		},
	}, nil

}

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
