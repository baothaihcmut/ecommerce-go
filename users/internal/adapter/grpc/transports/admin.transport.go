package transports

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/endpoints"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/request"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/response"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	"github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
)

type AdminServer struct {
	LoginHandler       grpc.Handler
	VerifyTokenHandler grpc.Handler
}

func NewAdminServer(e endpoints.AdminEndpoints, req request.AdminRequestMapper, res response.AdminResponseMapper) proto.AdminServiceServer {
	return &AdminServer{
		LoginHandler:       grpc.NewServer(e.LogIn, req.ToAdminLoginCommand, res.ToAdminLoginReponse),
		VerifyTokenHandler: grpc.NewServer(e.VerifyToken, req.ToAdminVerifyTokenCommand, res.ToAdminVerifyTokenResponse),
	}
}
func (a *AdminServer) LogIn(ctx context.Context, req *proto.AdminLoginRequest) (*proto.AdminLoginResponse, error) {
	_, resp, err := a.LoginHandler.ServeGRPC(ctx, req)
	if err != nil {
		return &proto.AdminLoginResponse{
			Data: nil,
			Status: &proto.Status{
				Code:    MapErrorToGrpcStatus(err).String(),
				Message: err.Error(),
				Details: []string{err.Error()},
			},
		}, nil
	}
	return &proto.AdminLoginResponse{
		Data: resp.(*proto.AdminLoginData),
		Status: &proto.Status{
			Code:    codes.OK.String(),
			Message: "Admin login success",
		},
	}, nil
}
func (a *AdminServer) VerifyToken(ctx context.Context, req *proto.AdminVerifyTokenRequest) (*proto.AdminVerifyTokenResponse, error) {
	_, resp, err := a.VerifyTokenHandler.ServeGRPC(ctx, req)
	if err != nil {
		return &proto.AdminVerifyTokenResponse{
			Data: nil,
			Status: &proto.Status{
				Code:    MapErrorToGrpcStatus(err).String(),
				Message: err.Error(),
				Details: []string{err.Error()},
			},
		}, nil
	}
	return &proto.AdminVerifyTokenResponse{
		Data: resp.(*proto.AdminVerifyTokenData),
		Status: &proto.Status{
			Code:    codes.OK.String(),
			Message: "Admin verify token success",
		},
	}, nil
}
