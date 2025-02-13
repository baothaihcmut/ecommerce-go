package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/admin/dtos/request"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/admin/dtos/response"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/admin/proto"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/models"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/google/uuid"
	grpcpool "github.com/processout/grpc-go-pool"
	"go.opentelemetry.io/otel/trace"
)

type AdminHandler interface {
	LogIn(context.Context, *request.AdminLoginRequestDTO) (*response.AdminLoginResponseDTO, error)
	VerifyToken(context.Context, string, bool) (*models.UserContext, error)
}

type AdminHandlerImpl struct {
	userAuthConnectionPool *grpcpool.Pool
	userConnectionPool     *grpcpool.Pool
	tracer                 trace.Tracer
}

func (a *AdminHandlerImpl) LogIn(ctx context.Context, dto *request.AdminLoginRequestDTO) (_ *response.AdminLoginResponseDTO, err error) {
	ctx, span := tracing.StartSpan(ctx, a.tracer, "Admin.LogIn: Call user service", nil)
	defer tracing.EndSpan(span, err, nil)
	conn, err := a.userConnectionPool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	s := proto.NewAdminServiceClient(conn)
	req := &proto.AdminLoginRequest{
		Email:    dto.Email,
		Password: dto.Password,
	}
	resp, err := s.LogIn(ctx, req)
	if err != nil {
		return nil, err
	}
	return &response.AdminLoginResponseDTO{
		AccessToken:  resp.Data.AccessToken,
		RefreshToken: resp.Data.RefreshToken,
	}, nil
}
func (a *AdminHandlerImpl) VerifyToken(ctx context.Context, token string, isRefreshToken bool) (_ *models.UserContext, err error) {
	ctx, span := tracing.StartSpan(ctx, a.tracer, "Auth.VerifyToken: Call user service", nil)
	defer tracing.EndSpan(span, err, nil)
	conn, err := a.userConnectionPool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	s := proto.NewAdminServiceClient(conn)
	req := proto.AdminVerifyTokenRequest{
		Token:          token,
		IsRefreshToken: isRefreshToken,
	}

	res, err := s.VerifyToken(ctx, &req)
	if err != nil {
		return nil, err
	}
	userId, err := uuid.Parse(res.Data.Id)
	if err != nil {
		return nil, err
	}
	return &models.UserContext{
		Id:   userId,
		Role: models.RoleAdmin,
	}, nil
}

func NewAdminHandler(
	userAuthConnectionPool *grpcpool.Pool,
	userConnectionPool *grpcpool.Pool,
	tracer trace.Tracer,
) AdminHandler {
	return &AdminHandlerImpl{
		userAuthConnectionPool: userAuthConnectionPool,
		userConnectionPool:     userConnectionPool,
		tracer:                 tracer,
	}
}
