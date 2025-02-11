package utils

import (
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/enums"
)

func MapRole(src enums.Role) *proto.Role {
	switch src {
	case enums.CUSTOMER:
		return proto.Role_CUSTOMER.Enum()
	case enums.SHOP_OWNER:
		return proto.Role_SHOP_OWNER.Enum()
	case enums.ADMIN:
		return proto.Role_ADMIN.Enum()
	default:
		return proto.Role_CUSTOMER.Enum()
	}
}
