package transports

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/endpoints"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/request"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/response"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
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
		return &proto.CreateUserResponse{
			Data: nil,
			Status: &proto.Status{
				Code:    codes.Internal.String(),
				Message: err.Error(),
				Details: []string{err.Error()},
				Success: false,
			},
		}, err
	}
	return &proto.CreateUserResponse{
		Data: resp.(*proto.UserData),
		Status: &proto.Status{
			Success: true,
			Message: "Create user successfully",
			Code:    codes.OK.String(),
			Details: []string{},
		},
	}, nil
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
