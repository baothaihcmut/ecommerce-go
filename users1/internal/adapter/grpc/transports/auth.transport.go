package transports

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/endpoints"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/request"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/response"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	gt "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	LoginHandler       gt.Handler
	SignUpHandler      gt.Handler
	VerifyTokenHandler gt.Handler
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
		VerifyTokenHandler: gt.NewServer(
			endpoints.VerifyToken,
			requestMapper.ToVerifyTokenCommand,
			responseMapper.ToVerifyTokenResult,
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
				Code:    MapErrorToGrpcStatus(err).String(),
				Details: []string{err.Error()},
			},
		}, nil
	}
	return &proto.LoginResponse{
		Data: resp.(*proto.LoginData),
		Status: &proto.Status{
			Message: "Login success",
			Code:    codes.OK.String(),
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
				Code:    MapErrorToGrpcStatus(err).String(),
			},
		}, status.Error(MapErrorToGrpcStatus(err), err.Error())
	}
	return &proto.SignUpResponse{
		Data: resp.(*proto.LoginData),
		Status: &proto.Status{
			Message: "Create user successfully",
			Details: []string{},
			Code:    codes.OK.String(),
		},
	}, nil
}

func (s *AuthServer) VerifyToken(ctx context.Context, req *proto.VerifyTokenRequest) (*proto.VerifyTokenResponse, error) {
	_, resp, err := s.VerifyTokenHandler.ServeGRPC(ctx, req)
	if err != nil {
		return &proto.VerifyTokenResponse{
			Data: nil,
			Status: &proto.Status{
				Code:    MapErrorToGrpcStatus(err).String(),
				Message: err.Error(),
				Details: []string{err.Error()},
			},
		}, status.Error(MapErrorToGrpcStatus(err), err.Error())
	}
	return &proto.VerifyTokenResponse{
		Data: resp.(*proto.VerifyTokenData),
		Status: &proto.Status{
			Message: "Verify token success",
			Details: []string{},
			Code:    codes.OK.String(),
		},
	}, nil
}
