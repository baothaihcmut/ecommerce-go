package request

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/commands"
)

type AuthRequestMapper interface {
	ToLoginCommand(context.Context, interface{}) (interface{}, error)
	ToSignUpCommand(context.Context, interface{}) (interface{}, error)
}

type AuthRequestMapperImpl struct {
}

func (m *AuthRequestMapperImpl) ToLoginCommand(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.LoginRequest)
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}
	return &commands.LoginCommand{
		Email:    *email,
		Password: req.Password,
	}, nil
}

func (m *AuthRequestMapperImpl) ToSignUpCommand(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.SignUpRequest)
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
	dest := &commands.SignUpCommand{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    req.Password,
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
func NewAuthRequestMapper() AuthRequestMapper {
	return &AuthRequestMapperImpl{}
}
