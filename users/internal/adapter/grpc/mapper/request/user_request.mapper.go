package request

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/commands"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
)

type UserRequestMapper interface {
	ToCreateUserCommand(_ context.Context, request interface{}) (interface{}, error)
}

type UserRequestMapperImpl struct {
}

func NewUserRequestMapper() UserRequestMapper {
	return &UserRequestMapperImpl{}
}

func (m *UserRequestMapperImpl) ToCreateUserCommand(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.CreateUserRequest)
	addreses := make([]*commands.Address, len(req.Addresses))
	for idx, addr := range req.Addresses {
		addreses[idx] = &commands.Address{
			Priority: int(addr.Priority),
			City:     addr.City,
			Town:     addr.Town,
			Street:   addr.Street,
			Province: addr.Province,
		}
	}
	dest := &commands.CreateUserCommand{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Role:        enums.Role(req.Role.String()),
		Addresses:   addreses,
	}
	if req.Role == proto.Role_CUSTOMER {
		dest.CustomerInfo = &commands.CustomerInfo{}
	} else {
		dest.ShopOwnerInfo = &commands.ShopOwnerInfo{
			BussinessLincese: req.ShopOwnerInfo.BussinessLincese,
		}

	}

	return dest, nil
}
