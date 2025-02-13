package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/utils"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/dtos/request"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/dtos/response"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/proto"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/models"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/google/uuid"
	grpcpool "github.com/processout/grpc-go-pool"
	"go.opentelemetry.io/otel/trace"
)

type AuthHandler interface {
	LogIn(context.Context, *request.LoginRequestDTO) (*response.LoginResponeDTO, error)
	SignUp(context.Context, *request.SignUpRequestDTO) (*response.LoginResponeDTO, error)
	VerifyToken(context.Context, string, bool) (*models.UserContext, error)
}

type UserAuthHandler struct {
	userAuthConnectionPool *grpcpool.Pool
	userConnectionPool     *grpcpool.Pool
	tracer                 trace.Tracer
}

func NewAuthHandler(
	userAuthConnectionPool *grpcpool.Pool,
	userConnectionPool *grpcpool.Pool,
	tracer trace.Tracer,
) AuthHandler {
	return &UserAuthHandler{
		tracer:                 tracer,
		userAuthConnectionPool: userAuthConnectionPool,
		userConnectionPool:     userConnectionPool,
	}
}

func mapAddress(addresses []request.Address) []*proto.Address {
	res := make([]*proto.Address, len(addresses))
	for idx, addr := range addresses {
		res[idx] = &proto.Address{
			Street:   addr.Street,
			City:     addr.City,
			Town:     addr.Town,
			Province: addr.Province,
			Priority: addr.Priority,
		}
	}
	return res
}
func (h *UserAuthHandler) LogIn(ctx context.Context, dto *request.LoginRequestDTO) (_ *response.LoginResponeDTO, err error) {
	ctx, span := tracing.StartSpan(ctx, h.tracer, "Auth.LogIn: Call user service", nil)
	defer tracing.EndSpan(span, err, nil)
	conn, err := h.userConnectionPool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	s := proto.NewAuthServiceClient(conn)
	req := &proto.LoginRequest{
		Email:    dto.Email,
		Password: dto.Password,
	}
	resp, err := s.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	return &response.LoginResponeDTO{
		AccessToken:  resp.Data.AccessToken,
		RefreshToken: resp.Data.RefreshToken,
	}, nil
}
func (h *UserAuthHandler) SignUp(ctx context.Context, dto *request.SignUpRequestDTO) (*response.LoginResponeDTO, error) {
	var err error
	ctx, span := tracing.StartSpan(ctx, h.tracer, "Auth.SignUp: Call user service", nil)
	defer tracing.EndSpan(span, err, nil)
	conn, err := h.userConnectionPool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	s := proto.NewAuthServiceClient(conn)
	req := &proto.SignUpRequest{
		Email:       dto.Email,
		Password:    dto.Password,
		PhoneNumber: dto.PhoneNumber,
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		Role:        utils.MapRole(dto.Role),
		Addresses:   mapAddress(dto.Addresses),
	}
	if dto.Role == models.RoleShopOwner {
		req.ShopOwnerInfo = &proto.ShopOwnerInfo{
			BussinessLincese: dto.ShopOwnerInfo.BusinessLicense,
		}
	}

	resp, err := s.SignUp(ctx, req)
	if err != nil {
		return nil, err
	}
	return &response.LoginResponeDTO{
		AccessToken:  resp.Data.AccessToken,
		RefreshToken: resp.Data.RefreshToken,
	}, nil
}

func (h *UserAuthHandler) VerifyToken(ctx context.Context, token string, isRefreshToken bool) (_ *models.UserContext, err error) {
	ctx, span := tracing.StartSpan(ctx, h.tracer, "Auth.VerifyToken: Call user service", nil)
	defer tracing.EndSpan(span, err, nil)
	conn, err := h.userConnectionPool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	s := proto.NewAuthServiceClient(conn)

	res, err := s.VerifyToken(ctx, &proto.VerifyTokenRequest{
		Token:          token,
		IsRefreshToken: isRefreshToken,
	})
	if err != nil {
		return nil, err
	}
	userId, err := uuid.Parse(res.Data.Id)
	if err != nil {
		return nil, err
	}
	return &models.UserContext{
		Id:   userId,
		Role: utils.MapRoleProto(*res.Data.Role),
	}, nil
}
