package handlers

import (
	"context"
	"fmt"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/utils"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/discovery"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/users/dtos"
	proto "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/users/proto/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserHandler interface {
	CreateUser(ctx context.Context, dto dtos.CreateUserDto) (*proto.CreateUserResponse, error)
}
type UserHandlerImpl struct {
	discoveryService discovery.DiscoveryService
}

func toPbCreateUser(dto *dtos.CreateUserDto) *proto.CreateUserRequest {
	addressPb := make([]*proto.Address, len(dto.Addresses))
	for idx, address := range dto.Addresses {
		addressPb[idx] = &proto.Address{
			Priority: int32(address.Priority),
			Street:   address.Street,
			Town:     address.Town,
			City:     address.City,
			Province: address.Province,
		}
	}
	return &proto.CreateUserRequest{
		Email:       dto.Email,
		PhoneNumber: dto.PhoneNumber,
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		Role:        utils.MapRole(dto.Role),
		Addresses:   addressPb,
		ShopOwnerInfo: &proto.ShopOwnerInfo{
			BussinessLincese: dto.ShopOwnerInfo.BusinessLicense,
		},
		CustomerInfo: &proto.CustomerInfo{},
	}
}

func NewUserHandler(discoveryService discovery.DiscoveryService) UserHandler {
	return &UserHandlerImpl{
		discoveryService: discoveryService,
	}
}
func (u *UserHandlerImpl) CreateUser(ctx context.Context, dto dtos.CreateUserDto) (*proto.CreateUserResponse, error) {
	//discovery service
	host, port, err := u.discoveryService.DiscoverService("user-service", "")
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	//call user service
	client := proto.NewUserServiceClient(conn)
	req := toPbCreateUser(&dto)
	resp, err := client.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil

}
