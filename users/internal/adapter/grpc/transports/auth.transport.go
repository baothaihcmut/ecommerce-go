package transports

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/endpoints"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/request"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/response"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	gt "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	LoginHandler  gt.Handler
	SignUpHandler gt.Handler
}

func NewAuthServer(
	endpoints endpoints.AuthEndpoints,
	requestMapper request.AuthRequestMapper,
	responseMapper response.AuthResponseMapper,
) proto.AuthServiceServer {
	return &AuthServer{
		LoginHandler: gt.NewServer(
			endpoints.Login,
			requestMapper.ToLoginCommand,
			responseMapper.ToLoginResult,
		),
		SignUpHandler: gt.NewServer(
			endpoints.SignUp,
			requestMapper.ToSignUpCommand,
			responseMapper.ToSignUpResult,
		),
	}
}

func (s *AuthServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	_, resp, err := s.LoginHandler.ServeGRPC(ctx, req)
	if err != nil {
		return &proto.LoginResponse{
			Data: nil,
			Status: &proto.Status{
				Message: err.Error(),
				Details: []string{err.Error()},
				Success: false,
			},
		}, status.Error(MapErrorToGrpcStatus(err), err.Error())
	}
	return &proto.LoginResponse{
		Data: resp.(*proto.LoginData),
		Status: &proto.Status{
			Success: true,
			Message: "Create user successfully",
			Details: []string{},
		},
	}, nil
}
func (s *AuthServer) SignUp(ctx context.Context, req *proto.SignUpRequest) (*proto.SignUpResponse, error) {
	_, resp, err := s.SignUpHandler.ServeGRPC(ctx, req)
	if err != nil {
		return &proto.SignUpResponse{
			Data: nil,
			Status: &proto.Status{
				Message: err.Error(),
				Details: []string{err.Error()},
				Success: false,
			},
		}, status.Error(MapErrorToGrpcStatus(err), err.Error())
	}
	return &proto.SignUpResponse{
		Data: resp.(*proto.LoginData),
		Status: &proto.Status{
			Success: true,
			Message: "Create user successfully",
			Details: []string{},
		},
	}, nil
}
