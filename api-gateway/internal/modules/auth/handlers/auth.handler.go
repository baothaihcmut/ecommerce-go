package handlers

import (
	"context"
	"fmt"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/grpc/interceptor"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/dtos/request"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/dtos/response"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/proto"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthHandler interface {
	LogIn(context.Context, *request.LoginRequestDTO) (*response.LoginResponeDTO, error)
	SignUp(context.Context, *request.SignUpRequestDTO) (*response.LoginResponeDTO, error)
}

type AuthHandlerImpl struct {
	discoveryService discovery.DiscoveryService
}

func NewAuthHandler(discoveryService discovery.DiscoveryService) AuthHandler {
	return &AuthHandlerImpl{
		discoveryService: discoveryService,
	}
}
func mapRole(src request.Role) proto.Role {
	switch src {
	case request.RoleCustomer:
		return proto.Role_CUSTOMER
	case request.RoleShopOwner:
		return proto.Role_SHOP_OWNER
	case request.RoleAdmin:
		return proto.Role_ADMIN
	}
	return proto.Role_CUSTOMER
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
func (h *AuthHandlerImpl) LogIn(ctx context.Context, dto *request.LoginRequestDTO) (*response.LoginResponeDTO, error) {
	host, port, err := h.discoveryService.DiscoverService("users-service", "")
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.ErrorHandlerClientInterceptor[proto.LoginData]()),
	)
	defer conn.Close()
	if err != nil {
		return nil, err
	}
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
func (h *AuthHandlerImpl) SignUp(ctx context.Context, dto *request.SignUpRequestDTO) (*response.LoginResponeDTO, error) {
	host, port, err := h.discoveryService.DiscoverService("users-service", "")
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.ErrorHandlerClientInterceptor[proto.LoginData]()),
	)
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	s := proto.NewAuthServiceClient(conn)
	req := &proto.SignUpRequest{
		Email:       dto.Email,
		Password:    dto.Password,
		PhoneNumber: dto.PhoneNumber,
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		Role:        mapRole(dto.Role),
		Addresses:   mapAddress(dto.Addresses),
	}
	if dto.Role == request.RoleShopOwner {
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
