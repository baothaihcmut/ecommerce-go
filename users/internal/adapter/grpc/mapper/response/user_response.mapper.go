package response

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
	"github.com/google/uuid"
)

type UserResponseMapper interface {
	ToCreateUserResponse(_ context.Context, response interface{}) (interface{}, error)
}

type UserResponseMapperImpl struct {
}

func (m *UserResponseMapperImpl) ToCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(*user.User)
	address := make([]*proto.Address, len(res.Address))
	for idx, addr := range res.Address {
		address[idx] = &proto.Address{
			Priority: int32(addr.Priority),
			City:     addr.City,
			Town:     addr.Town,
			Street:   addr.Street,
			Province: addr.Province,
		}
	}
	dest := &proto.UserData{
		Id:          uuid.UUID(res.Id).String(),
		FirstName:   res.FirstName,
		LastName:    res.LastName,
		Email:       string(res.Email),
		Role:        proto.Role(proto.Role_value[string(res.Role)]),
		PhoneNumber: string(res.PhoneNumber),
		Addresses:   address,
	}
	if res.Role == enums.CUSTOMER {
		dest.CustomerInfo = &proto.CustomerInfo{}
	} else {
		dest.ShopOwnerInfo = &proto.ShopOwnerInfo{
			BussinessLincese: res.ShopOwner.BussinessLicense,
		}
	}
	return dest, nil
}
func NewUserResponseMapper() UserResponseMapper {
	return &UserResponseMapperImpl{}
}
